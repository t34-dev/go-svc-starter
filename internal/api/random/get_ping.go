package random

import (
	"context"
	"github.com/t34-dev/go-svc-starter/pkg/api/random_v1"
)

func (s *ImplementedRandom) GetPing(context.Context, *random_v1.EmptyRequest) (*random_v1.PongResponse, error) {
	return &random_v1.PongResponse{
		Text: "pong",
	}, nil
}
