package auth_imp

import (
	"context"
	"github.com/t34-dev/go-svc-starter/pkg/api/auth_v1"
)

func (s *ImplementedAuth) Login(ctx context.Context, in *auth_v1.LoginRequest) (*auth_v1.LoginResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return &auth_v1.LoginResponse{
			AccessToken:  "AccessToken",
			RefreshToken: "RefreshToken",
		}, nil
	}
}
