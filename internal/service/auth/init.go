package auth_service

import (
	"context"
	"github.com/t34-dev/go-svc-starter/internal/service"
)

var jwtKey = []byte("your-secret-key")

var _ service.AuthService = &authService{}

type authService struct {
	opt service.Dependencies
}

func New(opt service.Dependencies) service.AuthService {
	return &authService{
		opt: opt,
	}
}

func (a authService) Check(ctx context.Context, path string) (bool, error) {
	return false, nil
}
