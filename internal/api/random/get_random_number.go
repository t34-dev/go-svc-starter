package random

import (
	"context"
	"github.com/t34-dev/go-svc-starter/pkg/api/random_v1"
	"math/rand"
)

func (s *ImplementedRandom) GetRandomNumber(ctx context.Context, _ *random_v1.EmptyRequest) (*random_v1.NumberResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return &random_v1.NumberResponse{Number: int32(rand.Intn(100))}, nil
	}
}
