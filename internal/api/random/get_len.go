package random

import (
	"context"
	"errors"
	"github.com/t34-dev/go-svc-starter/pkg/api/random_v1"
)

func (s *ImplementedRandom) GetLen(_ context.Context, req *random_v1.TxtRequest) (*random_v1.TxtResponse, error) {
	if req == nil {
		return nil, errors.New("GetLen: TxtRequest is empty")
	}

	return &random_v1.TxtResponse{
		Number: uint32(len(req.GetText())),
	}, nil
}
