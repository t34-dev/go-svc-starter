package pg_repository

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/t34-dev/go-svc-starter/internal/repository"
	commonRepository "github.com/t34-dev/go-svc-starter/internal/repository/pg/common"
	userRepository "github.com/t34-dev/go-svc-starter/internal/repository/pg/user"
)

func New(db *pgxpool.Pool) *repository.Repository {
	return &repository.Repository{
		Common: commonRepository.New(db),
		User:   userRepository.New(db),
	}
}
