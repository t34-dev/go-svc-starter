package common_service

import (
	"context"
	"github.com/t34-dev/go-svc-starter/internal/model"
	"github.com/t34-dev/go-svc-starter/internal/service"
	"github.com/t34-dev/go-utils/pkg/trace"
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
	ctx, finish := trace.TraceFunc(ctx, "commonService.GetDBTime", nil)
	defer finish()

	return s.opt.Repos.Common.GetTime(ctx)
}

func (s *commonService) GetTime(ctx context.Context) (time.Time, error) {
	ctx, finish := trace.TraceFunc(ctx, "commonService.GetTime", nil)
	defer finish()
	return time.Now(), nil
}

func (s *commonService) GetPost(ctx context.Context, id int64) (*model.Post, error) {
	ctx, finish := trace.TraceFunc(ctx, "commonService.GetPost", map[string]interface{}{"id": id})
	defer finish()

	post, err := s.opt.OtherService.GetPost(ctx, id)
	if err != nil {
		return nil, err
	}
	return post, err
}
