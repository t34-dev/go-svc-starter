package user_repository

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/t34-dev/go-svc-starter/internal/model"
	"github.com/t34-dev/go-svc-starter/internal/repository"
	"github.com/t34-dev/go-utils/pkg/db"
	"time"
)

const (
	userTable       = "users"
	idColumn        = "id"
	emailColumn     = "email"
	usernameColumn  = "username"
	passwordColumn  = "password"
	logoUrlColumn   = "logo_url"
	isBlockedColumn = "is_blocked"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
)

var _ repository.UserRepository = (*userRepository)(nil)

type userRepository struct {
	db      db.Client
	builder sq.StatementBuilderType
}

func New(db db.Client) repository.UserRepository {
	return &userRepository{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r userRepository) CreateUser(ctx context.Context, email, username, hashedPassword string) (uuid.UUID, error) {
	id := uuid.New()
	query, args, err := r.builder.Insert(userTable).
		Columns(idColumn, emailColumn, usernameColumn, passwordColumn, createdAtColumn, updatedAtColumn).
		Values(id, email, username, hashedPassword, time.Now(), time.Now()).
		ToSql()
	if err != nil {
		return uuid.Nil, err
	}

	q := db.Query{
		Name:     "user_repository.CreateUser",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	return id, err
}

func (r userRepository) GetUserByLogin(ctx context.Context, login string) (model.User, error) {
	query, args, err := r.builder.Select(idColumn, emailColumn, usernameColumn, passwordColumn).From(userTable).
		Where(sq.Or{
			sq.Eq{emailColumn: login},
			sq.Eq{usernameColumn: login},
		}).Limit(1).
		ToSql()
	if err != nil {
		return model.User{}, err
	}

	q := db.Query{
		Name:     "user_repository.GetUserByLogin",
		QueryRaw: query,
	}

	var user model.User
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&user.ID, &user.Email, &user.Username, &user.Password)
	return user, err
}

func (r userRepository) GetUserInfo(ctx context.Context, userId uuid.UUID) (model.User, error) {
	query, args, err := r.builder.Select(
		idColumn, emailColumn, usernameColumn, logoUrlColumn, isBlockedColumn, createdAtColumn, updatedAtColumn,
	).From(userTable).Where(sq.Eq{idColumn: userId}).
		ToSql()
	if err != nil {
		return model.User{}, err
	}

	q := db.Query{
		Name:     "user_repository.GetUserInfo",
		QueryRaw: query,
	}

	var user model.User
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(
		&user.ID, &user.Email, &user.Username, &user.LogoURL, &user.IsBlocked, &user.CreatedAt, &user.UpdatedAt,
	)
	return user, err
}
