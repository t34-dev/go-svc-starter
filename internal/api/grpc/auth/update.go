package auth

import (
	"context"
	"github.com/t34-dev/go-svc-starter/pkg/api/auth_v1"
)

func (s *ImplementedAuth) UpdateToken(ctx context.Context, in *auth_v1.UpdateTokenRequest) (*auth_v1.UpdateTokenResponse, error) {
	return &auth_v1.UpdateTokenResponse{
		AccessToken:  "AccessToken",
		RefreshToken: "RefreshToken",
	}, nil
}
