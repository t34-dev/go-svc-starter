package common

import (
	"github.com/t34-dev/go-svc-starter/internal/service"
	"github.com/t34-dev/go-svc-starter/pkg/api/common_v1"
)

var _ common_v1.CommonV1Server = &ImplementedCommon{}

type ImplementedCommon struct {
	common_v1.UnimplementedCommonV1Server
	service *service.Service
}

func NewImplementedCommon(service *service.Service) *ImplementedCommon {
	return &ImplementedCommon{service: service}
}
