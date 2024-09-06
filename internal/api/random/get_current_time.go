package random

import (
	"context"
	"github.com/t34-dev/go-svc-starter/pkg/api/random_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *ImplementedRandom) GetCurrentTime(ctx context.Context, _ *random_v1.EmptyRequest) (*random_v1.TimeResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return &random_v1.TimeResponse{Time: timestamppb.Now()}, nil
	}
}
