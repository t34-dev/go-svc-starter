package service

import (
	"context"
	"github.com/t34-dev/go-svc-starter/internal/repository"
	"time"
)

type Options struct {
	Repos *repository.Repository
}
type Service struct {
	Origin Options
	Common CommonService
}

//func New(repos Options) *Service {
//	srv := &Service{
//		Origin: repos,
//	}
//	srv.Common = common.New(srv)
//	return srv
//}

type CommonService interface {
	GetDBTime(ctx context.Context) (time.Time, error)
	GetTime(ctx context.Context) (time.Time, error)
}
type authInterface interface {
	// Login войти в систему
	Login(ctx context.Context, email, login, fingerPrint string, userAgent string) (string, error)
}
