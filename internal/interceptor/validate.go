package interceptor

import (
	"context"
	"google.golang.org/grpc"
)

type grpcValidate interface {
	Validate() error
}

func GrpcValidateInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	if val, ok := req.(grpcValidate); ok {
		if err := val.Validate(); err != nil {
			return nil, err
		}
	}
	return handler(ctx, req)
}
