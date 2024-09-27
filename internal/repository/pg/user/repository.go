package user_repository

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/t34-dev/go-svc-starter/internal/model"
	"github.com/t34-dev/go-svc-starter/internal/repository"
)

const (
	tableName = "note"

	idColumn        = "id"
	titleColumn     = "title"
	contentColumn   = "content"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
)

var _ repository.UserRepository = (*repo)(nil)

type repo struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) repository.UserRepository {
	return &repo{db: db}
}

func (r repo) GetAll(ctx context.Context, showIsBlock bool) ([]model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (r repo) Create(ctx context.Context, email, passwordHash, nickname string) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (r repo) UserById(ctx context.Context, id int) (*model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (r repo) UserByEmail(ctx context.Context, email string) (*model.User, error) {
	//TODO implement me
	panic("implement me")
}
