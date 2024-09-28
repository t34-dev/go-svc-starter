package common_repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/t34-dev/go-svc-starter/internal/repository"
)

var _ repository.CommonRepository = (*repo)(nil)

type repo struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func New(db *pgxpool.Pool) repository.CommonRepository {
	return &repo{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r repo) GetTime(ctx context.Context) (time.Time, error) {
	var dbTime time.Time

	query, args, err := r.builder.
		Select("NOW()").
		ToSql()

	if err != nil {
		return time.Time{}, err
	}

	err = r.db.QueryRow(ctx, query, args...).Scan(&dbTime)
	if err != nil {
		return time.Time{}, err
	}

	return dbTime, nil
}
