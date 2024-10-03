package othergrpcservice

import (
	"context"
	"github.com/t34-dev/go-svc-starter/internal/model"
)

type OtherGRPCService interface {
	GetPost(ctx context.Context, id int64) (*model.Post, error)
	GetPosts(ctx context.Context) ([]*model.Post, error)
}
