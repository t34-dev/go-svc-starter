package auth

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ImplementedAuth) Logout(ctx context.Context, in *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
