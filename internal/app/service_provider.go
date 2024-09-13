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
		// Строка подключения к базе данных
		databaseUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
			config.Pg().User(),
			config.Pg().Password(),
			config.Pg().Host(),
			config.Pg().Port(),
			config.Pg().DBName(),
			config.Pg().SSLMode(),
		)

		// Создаем конфигурацию для пула соединений
		dbConfig, err := pgxpool.ParseConfig(databaseUrl)
		if err != nil {
			logs.Fatal(fmt.Sprintf("Ошибка парсинга конфигурации: %v\n", err))
		}

		// Настраиваем параметры пула
		dbConfig.MinConns = config.Pg().MinConns()
		dbConfig.MaxConns = config.Pg().MaxConns()

		// Создаем пул соединений
		pool, err := pgxpool.ConnectConfig(ctx, dbConfig)
		if err != nil {
			logs.Fatal(fmt.Sprintf("Не удалось подключиться к базе данных: %v\n", err))
		}
		closer.Add(func() error {
			pool.Close()
			return nil
		})

		// Проверяем подключение
		if err := pool.Ping(ctx); err != nil {
			logs.Fatal(fmt.Sprintf("Не удалось выполнить пинг базы данных: %v\n", err))
		}
		logs.Debug("Успешно подключено к базе данных")
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
		srv := &service.Service{
			Origin: service.Options{
				Repos: s.Repos(ctx),
			},
		}
		srv.Common = commonSrv.New(srv)
		s.service = srv
	}

	return s.service
}
