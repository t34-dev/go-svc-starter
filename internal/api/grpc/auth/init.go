package auth

import (
	"github.com/t34-dev/go-svc-starter/internal/service"
	"github.com/t34-dev/go-svc-starter/pkg/api/auth_v1"
)

var _ auth_v1.AuthV1Server = &ImplementedAuth{}

type ImplementedAuth struct {
	auth_v1.UnimplementedAuthV1Server
	service *service.Service
}

func NewImplementedAuth(service *service.Service) *ImplementedAuth {
	return &ImplementedAuth{service: service}
}
