package common

import (
	"context"
	"github.com/t34-dev/go-svc-starter/pkg/api/common_v1"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *ImplementedCommon) GetTime(ctx context.Context, _ *emptypb.Empty) (*common_v1.TimeResponse, error) {
	return &common_v1.TimeResponse{Time: timestamppb.Now()}, nil
}
