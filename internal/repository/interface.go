package repository

import (
	"context"
	"github.com/t34-dev/go-svc-starter/internal/model"
	"time"
)

type Repository struct {
	Common CommonRepository
	User   UserRepository
}

type CommonRepository interface {
	GetTime(ctx context.Context) (time.Time, error)
}
type UserRepository interface {
	GetAll(ctx context.Context, showIsBlock bool) ([]model.User, error)
	Create(ctx context.Context, email, passwordHash, nickname string) (int, error)
	UserById(ctx context.Context, id int) (*model.User, error)
	UserByEmail(ctx context.Context, email string) (*model.User, error)
}
