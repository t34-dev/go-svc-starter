package app

import (
	"context"
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
	commonRepository "github.com/t34-dev/go-svc-starter/internal/repository/pg/common"
	role_repository "github.com/t34-dev/go-svc-starter/internal/repository/pg/role"
	deviceRepository "github.com/t34-dev/go-svc-starter/internal/repository/pg/session"
	userRepository "github.com/t34-dev/go-svc-starter/internal/repository/pg/user"
	"github.com/t34-dev/go-svc-starter/internal/service"
	accessService "github.com/t34-dev/go-svc-starter/internal/service/access"
	authService "github.com/t34-dev/go-svc-starter/internal/service/auth"
	commonService "github.com/t34-dev/go-svc-starter/internal/service/common"
	access_manager "github.com/t34-dev/go-svc-starter/pkg/access-manager"
	"github.com/t34-dev/go-svc-starter/pkg/api/common_v1"
	"github.com/t34-dev/go-utils/pkg/closer"
	"github.com/t34-dev/go-utils/pkg/db"
	"github.com/t34-dev/go-utils/pkg/db/pg"
	"github.com/t34-dev/go-utils/pkg/db/prettier"
	"github.com/t34-dev/go-utils/pkg/db/transaction"
	"github.com/t34-dev/go-utils/pkg/etcd"
	"github.com/t34-dev/go-utils/pkg/file"
	"github.com/t34-dev/go-utils/pkg/logs"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

const (
	modelPath  = "./model.conf"
	policyPath = "./policy.csv"
)

type serviceProvider struct {
	// grpc
	grpcImpl *grpcImpl.GrpcImpl

	// db
	dbPool    *pgxpool.Pool
	dbClient  db.Client
	txManager db.TxManager
	repos     *repository.Repository

	// service
	service *service.Service

	// grpc-clients
	clientOtherGrpc othergrpcservice.OtherGRPCService

	etcd          etcd.Client
	accessManager access_manager.AccessManager
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) DBPool(ctx context.Context) *pgxpool.Pool {
	if s.dbPool == nil {
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
			logs.Fatal(fmt.Sprintf("Error parsing configuration: %v", err))
		}

		// Configure pool parameters
		dbConfig.MinConns = config.Pg().MinConns()
		dbConfig.MaxConns = config.Pg().MaxConns()

		// Create connection pool
		pool, err := pgxpool.ConnectConfig(ctx, dbConfig)
		if err != nil {
			logs.Fatal(fmt.Sprintf("failed to connect to database: %v", err))
		}
		closer.Add(func() error {
			pool.Close()
			return nil
		})

		// Check connection
		if err := pool.Ping(ctx); err != nil {
			logs.Fatal(fmt.Sprintf("failed to ping database: %v", err))
		}
		return pool
	}
	return s.dbPool
}
func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		customLogger := pg.LogFunc(func(ctx context.Context, q db.Query, args ...interface{}) {
			prettyQuery := prettier.Pretty(q.QueryRaw, prettier.PlaceholderDollar, args...)
			logs.Info(fmt.Sprintf("%v", ctx), zap.String("SQL", q.Name), zap.String("query", prettyQuery))
		})
		cl, err := pg.New(s.DBPool(ctx), &customLogger)
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
func (s *serviceProvider) ETCD(_ context.Context) etcd.Client {
	if s.etcd == nil {
		etcdConfig := clientv3.Config{
			Endpoints: []string{"localhost:2378"},
		}
		cli, err := etcd.NewClient(etcdConfig, nil)
		if err != nil {
			logs.Fatal("failed to create etcd client", zap.Error(err))
		}
		closer.Add(cli.Close)
		s.etcd = cli
	}

	return s.etcd
}
func (s *serviceProvider) AccessManager(ctx context.Context) access_manager.AccessManager {
	if s.accessManager == nil {
		accessManager := access_manager.NewAccessManager(s.ETCD(ctx))
		var err error

		if err = accessManager.UpdateConfigsFromFiles(modelPath, policyPath); err != nil {
			logs.Fatal("failed to update config from files", zap.Error(err))
		}

		// WATCHER
		callback := file.FileChangeCallback(func(path string, newData []byte, err error) {
			if err != nil {
				logs.Error(fmt.Sprintf("error for file %s: %v", path), zap.Error(err))
				return
			}
			logs.Warn(fmt.Sprintf("config ROLE file %s updated locally", path), zap.String("newConfig", string(newData)))

			if err = accessManager.UpdateConfigsFromFiles(modelPath, policyPath); err != nil {
				logs.Error("failed to update config from files", zap.Error(err))
				return
			}

			if err = accessManager.UpdateEtcdStore(ctx); err != nil {
				logs.Error("failed to update etcd store", zap.Error(err))
				return
			}
		})
		watcher, err := file.NewWatcher(callback)
		if err != nil {
			logs.Fatal("failed creating watcher", zap.Error(err))
		}

		err = watcher.WatchFiles([]string{
			modelPath,
			policyPath,
		})

		if err = accessManager.UpdateConfigsFromFiles(modelPath, policyPath); err != nil {
			logs.Fatal("failed to update config from files", zap.Error(err))
		}

		if err = accessManager.UpdateEtcdStore(ctx); err != nil {
			logs.Fatal("failed to update etcd store", zap.Error(err))
		}

		// попытка обновить конфигурацию из etcd хранилища
		if err = accessManager.UpdateConfigsFromEtcd(ctx); err != nil {
			logs.Fatal("failed to update ETCD config", zap.Error(err))
		}
		// обновить конфигурацию только если конфигурация поменялась в хранилище
		err = accessManager.WatchConfig(ctx, func(err2 error, key string, newValue []byte) {
			if err2 != nil {
				logs.Error("failed to update watch config", zap.Error(err2))
				return
			}
			logs.Warn("updated watch config", zap.String("key", key), zap.String("newValue", string(newValue)))
		})
		if err != nil {
			logs.Fatal("failed to update watch config", zap.Error(err))
		}
		// при завершении перестать отлеживать изменения
		closer.Add(accessManager.StopWatchConfig)
		s.accessManager = accessManager
	}

	return s.accessManager
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}
func (s *serviceProvider) Repos(ctx context.Context) *repository.Repository {
	if s.repos == nil {
		dbClient := s.DBClient(ctx)
		s.repos = &repository.Repository{
			Common:  commonRepository.New(dbClient),
			User:    userRepository.New(dbClient),
			Session: deviceRepository.New(dbClient),
			Role:    role_repository.New(dbClient),
		}
	}
	return s.repos
}
func (s *serviceProvider) OtherGrpc(_ context.Context) othergrpcservice.OtherGRPCService {
	if s.clientOtherGrpc == nil {
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

		s.clientOtherGrpc = othergrpcservice_impl.New(commonSrv)
	}
	return s.clientOtherGrpc
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
		deps := service.NewDeps(srv, *s.Repos(ctx), s.OtherGrpc(ctx), s.AccessManager(ctx))
		srv.Common = commonService.New(deps)
		srv.Auth = authService.New(deps)
		srv.Access = accessService.New(deps)
		s.service = &srv
	}

	return s.service
}
