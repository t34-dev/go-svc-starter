package common

import (
	"context"
	"github.com/t34-dev/go-svc-starter/pkg/api/common_v1"
	"github.com/t34-dev/go-utils/pkg/sys"
	"github.com/t34-dev/go-utils/pkg/sys/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *ImplementedCommon) GetTime(ctx context.Context, req *common_v1.TimeRequest) (*common_v1.TimeResponse, error) {
	if req.GetError() {
		return nil, sys.NewCommonError("Error message!", codes.Internal)
	}

	tt, err := s.service.Common.GetTime(ctx)
	if err != nil {
		return nil, sys.NewCommonError("ERR", codes.Internal)
	}
	return &common_v1.TimeResponse{Time: timestamppb.New(tt)}, nil
}
