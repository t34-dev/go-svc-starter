package random

import (
	"context"
	"github.com/t34-dev/go-svc-starter/pkg/api/random_v1"
)

func (s *ImplementedRandom) GetPerson(context.Context, *random_v1.EmptyRequest) (*random_v1.Person, error) {
	children := []string{
		"zak1", "zak2", "zak3",
	}
	return &random_v1.Person{
		Name:     "T34",
		Age:      22,
		Children: children,
		Parent: &random_v1.Parent{
			Name: "parent",
		},
	}, nil
}
