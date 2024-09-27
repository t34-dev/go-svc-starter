package common_service

import (
	"context"
	"github.com/t34-dev/go-svc-starter/internal/service"
	"time"
)

var _ service.CommonService = &commonService{}

type commonService struct {
	opt service.Dependencies
}

func New(opt service.Dependencies) service.CommonService {
	return &commonService{
		opt: opt,
	}
}

func (s *commonService) GetDBTime(ctx context.Context) (time.Time, error) {
	return s.opt.Repos.Common.GetTime(ctx)
}

func (s *commonService) GetTime(_ context.Context) (time.Time, error) {
	return time.Now(), nil
}
