package interceptor

import (
	"context"
	"google.golang.org/grpc"
)

type validate interface {
	Validate() error
}

func ValidateInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	if val, ok := req.(validate); ok {
		if err := val.Validate(); err != nil {
			return nil, err
		}
	}
	return handler(ctx, req)
}
