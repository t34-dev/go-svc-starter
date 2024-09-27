package service

import (
	"context"
	"github.com/t34-dev/go-svc-starter/internal/repository"
	"time"
)

type Dependencies struct {
	Service Service
	Repos   repository.Repository
}

func NewDeps(service Service, repos repository.Repository) Dependencies {
	return Dependencies{
		Service: service,
		Repos:   repos,
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
}
type AccessService interface {
	Check(ctx context.Context, path string) (bool, error)
}
type AuthService interface {
	Check(ctx context.Context, path string) (bool, error)
}
