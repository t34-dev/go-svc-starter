package grpc_impl

import (
	"github.com/t34-dev/go-svc-starter/internal/api/grpc/access"
	"github.com/t34-dev/go-svc-starter/internal/api/grpc/auth"
	"github.com/t34-dev/go-svc-starter/internal/api/grpc/common"
)

type GrpcImpl struct {
	Common *common.ImplementedCommon
	Access *access.ImplementedAccess
	Auth   *auth.ImplementedAuth
}
