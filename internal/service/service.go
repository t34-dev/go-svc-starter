package service

import (
	"context"
	othergrpcservice "github.com/t34-dev/go-svc-starter/internal/client/other_grpc_service"
	"github.com/t34-dev/go-svc-starter/internal/model"
	"github.com/t34-dev/go-svc-starter/internal/repository"
	"time"
)

type Dependencies struct {
	Service      Service
	Repos        repository.Repository
	OtherService othergrpcservice.OtherGRPCService
}

func NewDeps(service Service, repos repository.Repository, otherService othergrpcservice.OtherGRPCService) Dependencies {
	return Dependencies{
		Service:      service,
		Repos:        repos,
		OtherService: otherService,
	}
}

type Service struct {
	Common CommonService
	Access AccessService
	Auth   AuthService
}

type CommonService interface {
	GetDBTime(ctx context.Context) (time.Time, error)
	GetTime(ctx context.Context) (time.Time, error)
	GetPost(ctx context.Context, id int64) (*model.Post, error)
}
type AccessService interface {
	Check(ctx context.Context, path string) (bool, error)
}
type AuthService interface {
	Check(ctx context.Context, path string) (bool, error)
}
