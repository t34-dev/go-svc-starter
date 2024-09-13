package access

import (
	"github.com/t34-dev/go-svc-starter/internal/service"
	"github.com/t34-dev/go-svc-starter/pkg/api/access_v1"
)

var _ access_v1.AccessV1Server = &ImplementedAccess{}

type ImplementedAccess struct {
	access_v1.UnimplementedAccessV1Server
	service *service.Service
}

func NewImplementedAccess(service *service.Service) *ImplementedAccess {
	return &ImplementedAccess{service: service}
}
