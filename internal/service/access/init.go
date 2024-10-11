package access_service

import (
	"context"
	"github.com/t34-dev/go-svc-starter/internal/service"
)

var _ service.AccessService = &accessService{}

type accessService struct {
	opt service.Options
}

func New(opt service.Options) service.AccessService {
	return &accessService{
		opt: opt,
	}
}

func (a accessService) Check(ctx context.Context, path string) (bool, error) {
	return false, nil
}
