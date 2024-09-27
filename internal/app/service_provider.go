package app

import (
	"context"
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
	accessSrv "github.com/t34-dev/go-svc-starter/internal/service/access"
	authSrv "github.com/t34-dev/go-svc-starter/internal/service/auth"
	commonSrv "github.com/t34-dev/go-svc-starter/internal/service/common"
	"github.com/t34-dev/go-utils/pkg/closer"
	"github.com/t34-dev/go-utils/pkg/logs"
)

type serviceProvider struct {
	db      *pgxpool.Pool
	repos   *repository.Repository
	service *service.Service

	// grpc
	grpcImpl *grpcImpl.GrpcImpl
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) DB(ctx context.Context) *pgxpool.Pool {
	if s.db == nil {
		// Database connection string
		databaseUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
			config.Pg().User(),
			config.Pg().Password(),
			config.Pg().Host(),
			config.Pg().Port(),
			config.Pg().DBName(),
			config.Pg().SSLMode(),
		)

		// Create configuration for the connection pool
		dbConfig, err := pgxpool.ParseConfig(databaseUrl)
		if err != nil {
			logs.Fatal(fmt.Sprintf("Error parsing configuration: %v\n", err))
		}

		// Configure pool parameters
		dbConfig.MinConns = config.Pg().MinConns()
		dbConfig.MaxConns = config.Pg().MaxConns()

		// Create connection pool
		pool, err := pgxpool.ConnectConfig(ctx, dbConfig)
		if err != nil {
			logs.Fatal(fmt.Sprintf("Failed to connect to database: %v\n", err))
		}
		closer.Add(func() error {
			pool.Close()
			return nil
		})

		// Check connection
		if err := pool.Ping(ctx); err != nil {
			logs.Fatal(fmt.Sprintf("Failed to ping database: %v\n", err))
		}
		logs.Debug("Successfully connected to database")
		s.db = pool
	}
	return s.db
}

func (s *serviceProvider) Repos(ctx context.Context) *repository.Repository {
	if s.repos == nil {
		s.repos = pgRepos.New(s.DB(ctx))
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
		srv.Common = commonSrv.New(deps)
		srv.Auth = authSrv.New(deps)
		srv.Access = accessSrv.New(deps)
		s.service = &srv
	}

	return s.service
}
