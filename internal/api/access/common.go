package access

import (
	"github.com/t34-dev/go-svc-starter/pkg/api/access_v1"
)

var _ access_v1.AccessV1Server = &ImplementedAccess{}

type ImplementedAccess struct {
	access_v1.UnimplementedAccessV1Server
}

func NewImplementedAccess() *ImplementedAccess {
	return &ImplementedAccess{}
}
