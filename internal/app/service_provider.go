package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/opentracing/opentracing-go"
	grpcpool "github.com/t34-dev/go-grpc-pool"
	grpcImpl "github.com/t34-dev/go-svc-starter/internal/api/grpc"
	"github.com/t34-dev/go-svc-starter/internal/api/grpc/access"
	"github.com/t34-dev/go-svc-starter/internal/api/grpc/auth"
	"github.com/t34-dev/go-svc-starter/internal/api/grpc/common"
	othergrpcservice "github.com/t34-dev/go-svc-starter/internal/client/other_grpc_service"
	othergrpcservice_impl "github.com/t34-dev/go-svc-starter/internal/client/other_grpc_service/impl"
	"github.com/t34-dev/go-svc-starter/internal/config"
	"github.com/t34-dev/go-svc-starter/internal/repository"
	pgRepos "github.com/t34-dev/go-svc-starter/internal/repository/pg"
	"github.com/t34-dev/go-svc-starter/internal/service"
	accessService "github.com/t34-dev/go-svc-starter/internal/service/access"
	authService "github.com/t34-dev/go-svc-starter/internal/service/auth"
	commonService "github.com/t34-dev/go-svc-starter/internal/service/common"
	"github.com/t34-dev/go-svc-starter/pkg/api/common_v1"
	"github.com/t34-dev/go-utils/pkg/closer"
	"github.com/t34-dev/go-utils/pkg/db"
	"github.com/t34-dev/go-utils/pkg/db/pg"
	"github.com/t34-dev/go-utils/pkg/db/prettier"
	"github.com/t34-dev/go-utils/pkg/db/transaction"
	"github.com/t34-dev/go-utils/pkg/logs"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type serviceProvider struct {
	dbClient  db.Client
	txManager db.TxManager
	repos     *repository.Repository
	service   *service.Service

	// grpc
	grpcImpl *grpcImpl.GrpcImpl

	// grpc
	otherGrpc othergrpcservice.OtherGRPCService
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		connector := func(ctx context.Context) (*pgxpool.Pool, error) {
			dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
				config.Pg().User(),
				config.Pg().Password(),
				config.Pg().Host(),
				config.Pg().Port(),
				config.Pg().DBName(),
				config.Pg().SSLMode(),
			)
			// Create configuration for the connection pool
			dbConfig, err := pgxpool.ParseConfig(dsn)
			if err != nil {
				logs.Fatal(fmt.Sprintf("Error parsing configuration: %v\n", err))
			}

			// Configure pool parameters
			dbConfig.MinConns = config.Pg().MinConns()
			dbConfig.MaxConns = config.Pg().MaxConns()

			// Create connection pool
			pool, err := pgxpool.ConnectConfig(ctx, dbConfig)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Failed to connect to database: %v\n", err))
			}
			closer.Add(func() error {
				pool.Close()
				return nil
			})

			// Check connection
			if err := pool.Ping(ctx); err != nil {
				return nil, errors.New(fmt.Sprintf("Failed to ping database: %v\n", err))
			}
			return pool, nil
		}
		customLogger := pg.LogFunc(func(ctx context.Context, q db.Query, args ...interface{}) {
			prettyQuery := prettier.Pretty(q.QueryRaw, prettier.PlaceholderDollar, args...)
			logs.Info(fmt.Sprintf("%v", ctx), zap.String("SQL", q.Name), zap.String("query", prettyQuery))
		})
		cl, err := pg.New(ctx, connector, &customLogger)
		if err != nil {
			logs.Fatal("failed to create db client", zap.Error(err))
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			logs.Fatal("failed to ping db", zap.Error(err))
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}
func (s *serviceProvider) Repos(ctx context.Context) *repository.Repository {
	if s.repos == nil {
		s.repos = pgRepos.New(s.DBClient(ctx))
	}
	return s.repos
}
func (s *serviceProvider) OtherGrpc(_ context.Context) othergrpcservice.OtherGRPCService {
	if s.otherGrpc == nil {
		creds, err := credentials.NewClientTLSFromFile("cert/service.pem", "")
		if err != nil {
			logs.Fatal("failed to load client TLS credentials:", zap.Error(err))
		}
		opts := []grpc.DialOption{
			grpc.WithBlock(),
			grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
		}

		if config.App().IsTSL() {
			opts = append(opts, grpc.WithTransportCredentials(creds))
		} else {
			opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
		}

		// Define a factory function to create gRPC connections
		factory := func() (*grpc.ClientConn, error) {
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()

			return grpc.DialContext(ctx, config.MS().OtherGrpcAddress(), opts...)
		}

		// Create a new connection pool
		grpcPool, err := grpcpool.NewPool(factory, grpcpool.PoolOptions{
			MinConn: 2,
			MaxConn: 30,
		})
		if err != nil {
			logs.Fatal(fmt.Sprintf("failed to connect to GRPC: %s", config.MS().OtherGrpcAddress()), zap.Error(err))
		}
		closer.Add(func() error {
			grpcPool.Close()
			return nil
		})
		conn, err := grpcPool.Get()
		if err != nil {
			logs.Fatal("did not connect:", zap.Error(err))
		}

		commonSrv := common_v1.NewCommonV1Client(conn.GetConn())

		s.otherGrpc = othergrpcservice_impl.New(commonSrv)
	}
	return s.otherGrpc
}

func (s *serviceProvider) GrpcImpl(ctx context.Context) *grpcImpl.GrpcImpl {
	if s.grpcImpl == nil {
		s.grpcImpl = &grpcImpl.GrpcImpl{
			Access: access.NewImplementedAccess(s.Service(ctx)),
			Auth:   auth.NewImplementedAuth(s.Service(ctx)),
			Common: common.NewImplementedCommon(s.Service(ctx)),
		}
	}

	return s.grpcImpl
}

func (s *serviceProvider) Service(ctx context.Context) *service.Service {
	if s.service == nil {
		srv := service.Service{}
		deps := service.NewDeps(srv, *s.Repos(ctx), s.OtherGrpc(ctx))
		srv.Common = commonService.New(deps)
		srv.Auth = authService.New(deps)
		srv.Access = accessService.New(deps)
		s.service = &srv
	}

	return s.service
}
