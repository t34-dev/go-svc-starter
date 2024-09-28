package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	grpcImpl "github.com/t34-dev/go-svc-starter/internal/api/grpc"
	"github.com/t34-dev/go-svc-starter/internal/api/grpc/access"
	"github.com/t34-dev/go-svc-starter/internal/api/grpc/auth"
	"github.com/t34-dev/go-svc-starter/internal/api/grpc/common"
	"github.com/t34-dev/go-svc-starter/internal/config"
	"github.com/t34-dev/go-svc-starter/internal/repository"
	pgRepos "github.com/t34-dev/go-svc-starter/internal/repository/pg"
	"github.com/t34-dev/go-svc-starter/internal/service"
	accessService "github.com/t34-dev/go-svc-starter/internal/service/access"
	authService "github.com/t34-dev/go-svc-starter/internal/service/auth"
	commonService "github.com/t34-dev/go-svc-starter/internal/service/common"
	"github.com/t34-dev/go-utils/pkg/closer"
	"github.com/t34-dev/go-utils/pkg/db"
	"github.com/t34-dev/go-utils/pkg/db/pg"
	"github.com/t34-dev/go-utils/pkg/db/prettier"
	"github.com/t34-dev/go-utils/pkg/db/transaction"
	"github.com/t34-dev/go-utils/pkg/logs"
	"go.uber.org/zap"
)

type serviceProvider struct {
	dbClient  db.Client
	txManager db.TxManager
	repos     *repository.Repository
	service   *service.Service

	// grpc
	grpcImpl *grpcImpl.GrpcImpl
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
		deps := service.NewDeps(srv, *s.Repos(ctx))
		srv.Common = commonService.New(deps)
		srv.Auth = authService.New(deps)
		srv.Access = accessService.New(deps)
		s.service = &srv
	}

	return s.service
}
