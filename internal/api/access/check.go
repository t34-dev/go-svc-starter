package access

import (
	"context"
	"fmt"
	"github.com/t34-dev/go-svc-starter/pkg/api/access_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ImplementedAccess) Check(ctx context.Context, in *access_v1.CheckRequest) (*emptypb.Empty, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		fmt.Println("Endpoin:", in.GetEndpointAddress())
		return &emptypb.Empty{}, nil
	}
}
