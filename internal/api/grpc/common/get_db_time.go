package common

import (
	"context"
	"github.com/t34-dev/go-svc-starter/pkg/api/common_v1"
	"github.com/t34-dev/go-utils/pkg/sys"
	"github.com/t34-dev/go-utils/pkg/sys/codes"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *ImplementedCommon) GetDBTime(ctx context.Context, _ *emptypb.Empty) (*common_v1.TimeResponse, error) {
	tt, err := s.service.Common.GetDBTime(ctx)
	if err != nil {
		return nil, sys.NewError(err.Error(), codes.Internal)
	}
	return &common_v1.TimeResponse{Time: timestamppb.New(tt)}, err
}
