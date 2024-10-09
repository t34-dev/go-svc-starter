package common_repository

import (
	"context"
	"github.com/t34-dev/go-utils/pkg/db"
	"github.com/t34-dev/go-utils/pkg/trace"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/t34-dev/go-svc-starter/internal/repository"
)

var _ repository.CommonRepository = (*commonRepository)(nil)

type commonRepository struct {
	db      db.Client
	builder sq.StatementBuilderType
}

func New(db db.Client) repository.CommonRepository {
	return &commonRepository{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r commonRepository) GetTime(ctx context.Context) (time.Time, error) {
	ctx, finish := trace.TraceFunc(ctx, "commonRepository.GetTime", nil)
	defer finish()

	var dbTime time.Time

	query, args, err := r.builder.
		Select("NOW()").
		ToSql()

	if err != nil {
		return time.Time{}, err
	}

	q := db.Query{
		Name:     "common_repository.GetTime",
		QueryRaw: query,
	}

	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&dbTime)
	if err != nil {
		return time.Time{}, err
	}

	return dbTime, nil
}
