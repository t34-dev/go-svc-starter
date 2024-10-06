package repository

import (
	"context"
	"github.com/t34-dev/go-svc-starter/internal/model"
	"time"
)

type Repository struct {
	Common  CommonRepository
	User    UserRepository
	Session SessionRepository
}

type CommonRepository interface {
	GetTime(ctx context.Context) (time.Time, error)
}
type UserRepository interface {
	CreateUser(ctx context.Context, email, username, hashedPassword string) (int64, error)
	GetUserByLogin(ctx context.Context, login string) (model.User, error)
	GetUserInfo(ctx context.Context, userID int64) (model.User, error)
}
type SessionRepository interface {
	CreateSession(ctx context.Context, userID int64, sessionKey, sessionName, refreshToken string) error
	UpsertSession(ctx context.Context, userID int64, sessionKey, sessionName, refreshToken string) error
	DeleteSession(ctx context.Context, userID int64, refreshToken string) error
	GetSessionByRefreshToken(ctx context.Context, refreshToken, sessionKey string) (model.Session, error)
	UpdateSession(ctx context.Context, userID int64, sessionKey, sessionName, refreshToken string) error
	GetActiveSessions(ctx context.Context, userID int64) ([]model.Session, error)
	GetCurrentSession(ctx context.Context, userID int64) (int64, error)
	CleanupInactiveSessions(ctx context.Context) error
}
