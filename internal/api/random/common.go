package random

import "github.com/t34-dev/go-svc-starter/pkg/api/random_v1"

var _ random_v1.RandomServiceServer = &ImplementedRandom{}

type ImplementedRandom struct {
	random_v1.UnimplementedRandomServiceServer
}

func NewImplementedRandom() *ImplementedRandom {
	return &ImplementedRandom{}
}
