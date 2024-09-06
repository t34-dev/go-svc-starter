package random

import (
	"context"
	"github.com/t34-dev/go-svc-starter/pkg/api/random_v1"
	"math/rand"
)

func (s *ImplementedRandom) GetRandomQuote(ctx context.Context, _ *random_v1.EmptyRequest) (*random_v1.QuoteResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		quotes := []string{
			"The only way to do great work is to love what you do.",
			"Life is what happens when you're busy making other plans.",
			"The future belongs to those who believe in the beauty of their dreams.",
		}
		return &random_v1.QuoteResponse{Quote: quotes[rand.Intn(len(quotes))]}, nil
	}
}
