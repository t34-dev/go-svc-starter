package app

import "github.com/t34-dev/go-svc-starter/internal/config"

type serviceProvider struct {
	appConfig  config.AppConfig
	grpcConfig config.GrpcConfig
	httpConfig config.HttpConfig
	pgConfig   config.PgConfig
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) APPConfig() config.AppConfig {
	if s.appConfig == nil {
		s.appConfig = config.App()
	}

	return s.appConfig
}

func (s *serviceProvider) PGConfig() config.PgConfig {
	if s.pgConfig == nil {
		s.pgConfig = config.Pg()
	}

	return s.pgConfig
}

func (s *serviceProvider) GRPCConfig() config.GrpcConfig {
	if s.grpcConfig == nil {
		s.grpcConfig = config.Grpc()
	}

	return s.grpcConfig
}

func (s *serviceProvider) HTTPConfig() config.HttpConfig {
	if s.httpConfig == nil {
		s.httpConfig = config.Http()
	}

	return s.httpConfig
}
