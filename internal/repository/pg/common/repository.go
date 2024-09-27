package common_repository

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/t34-dev/go-svc-starter/internal/repository"
	"time"
)

var _ repository.CommonRepository = (*repo)(nil)

type repo struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) repository.CommonRepository {
	return &repo{db: db}
}

func (r repo) GetTime(ctx context.Context) (time.Time, error) {
	var dbTime time.Time
	err := r.db.QueryRow(ctx, "SELECT NOW()").Scan(&dbTime)
	if err != nil {
		return time.Time{}, err
	}
	return dbTime, nil
}
