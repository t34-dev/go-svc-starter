package pg_repository

import (
	"github.com/t34-dev/go-svc-starter/internal/repository"
	commonRepository "github.com/t34-dev/go-svc-starter/internal/repository/pg/common"
	userRepository "github.com/t34-dev/go-svc-starter/internal/repository/pg/user"
	"github.com/t34-dev/go-utils/pkg/db"
)

func New(db db.Client) *repository.Repository {
	return &repository.Repository{
		Common: commonRepository.New(db),
		User:   userRepository.New(db),
	}
}
