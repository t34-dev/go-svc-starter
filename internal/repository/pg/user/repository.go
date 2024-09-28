package user_repository

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/t34-dev/go-svc-starter/internal/model"
	"github.com/t34-dev/go-svc-starter/internal/repository"
	"github.com/t34-dev/go-utils/pkg/db"
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
	db      db.Client
	builder sq.StatementBuilderType
}

func New(db db.Client) repository.UserRepository {
	return &repo{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
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
