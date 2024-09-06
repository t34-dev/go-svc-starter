package auth

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ImplementedAuth) Logout(ctx context.Context, in *emptypb.Empty) (*emptypb.Empty, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return &emptypb.Empty{}, nil
	}
}
