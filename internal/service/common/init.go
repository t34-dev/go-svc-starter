package common_srv

import (
	"context"
	"github.com/t34-dev/go-svc-starter/internal/service"
	"time"
)

var _ service.CommonService = &commonService{}

type commonService struct {
	origin *service.Service
}

func New(service *service.Service) service.CommonService {
	return &commonService{
		origin: service,
	}
}

func (s *commonService) GetDBTime(ctx context.Context) (time.Time, error) {
	return time.Now(), nil
}

func (s *commonService) GetTime(ctx context.Context) (time.Time, error) {
	return time.Now(), nil
}
