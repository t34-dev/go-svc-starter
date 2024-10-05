package user_repository

import (
	"context"
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/t34-dev/go-svc-starter/internal/model"
	"github.com/t34-dev/go-svc-starter/internal/repository"
	"time"
)

const (
	userTable           = "users"
	userIDColumn        = "id"
	userEmailColumn     = "email"
	userUsernameColumn  = "username"
	userPasswordColumn  = "password"
	userCreatedAtColumn = "created_at"
	userUpdatedAtColumn = "updated_at"
)

var _ repository.UserRepository = (*userRepository)(nil)

type userRepository struct {
	db      *sql.DB
	builder sq.StatementBuilderType
}

func New(db *sql.DB) repository.UserRepository {
	return &userRepository{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r userRepository) CreateUser(ctx context.Context, email, username, hashedPassword string) (int64, error) {
	var id int64
	err := r.builder.Insert(userTable).
		Columns(userEmailColumn, userUsernameColumn, userPasswordColumn, userCreatedAtColumn, userUpdatedAtColumn).
		Values(email, username, hashedPassword, time.Now(), time.Now()).
		Suffix("RETURNING " + userIDColumn).
		RunWith(r.db).
		QueryRow().
		Scan(&id)
	return id, err
}

func (r userRepository) GetUserByLogin(ctx context.Context, login string) (model.User, error) {
	var user model.User
	err := r.builder.Select(userIDColumn, userPasswordColumn).From(userTable).
		Where(sq.Or{
			sq.Eq{userEmailColumn: login},
			sq.Eq{userUsernameColumn: login},
		}).Limit(1).
		RunWith(r.db).
		QueryRow().
		Scan(&user.ID, &user.Password)
	return user, err
}

func (r userRepository) GetUserInfo(ctx context.Context, userId int64) (model.User, error) {
	var user model.User
	err := r.builder.Select(
		userIDColumn, userEmailColumn, userUsernameColumn, userCreatedAtColumn, userUpdatedAtColumn,
	).From(userTable).Where(sq.Eq{userIDColumn: userId}).
		RunWith(r.db).
		QueryRow().
		Scan(&user.ID, &user.Email, &user.Username, &user.CreatedAt, &user.UpdatedAt)
	return user, err
}
