package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/t34-dev/go-svc-starter/internal/model"
	"time"
)

type Repository struct {
	Common  CommonRepository
	User    UserRepository
	Session SessionRepository
	Role    RoleRepository
}

type CommonRepository interface {
	GetTime(ctx context.Context) (time.Time, error)
}
type UserRepository interface {
	CreateUser(ctx context.Context, email, username, hashedPassword string) (uuid.UUID, error)
	GetUserByLogin(ctx context.Context, login string) (model.User, error)
	GetUserInfo(ctx context.Context, userId uuid.UUID) (model.User, error)
}
type SessionRepository interface {
	CreateSession(ctx context.Context, userID uuid.UUID, deviceKey, deviceName string) (uuid.UUID, error)
	UpsertSession(ctx context.Context, userID uuid.UUID, deviceKey, deviceName string) (uuid.UUID, error)
	DeleteSession(ctx context.Context, sessionID uuid.UUID) error
	GetSessionByID(ctx context.Context, sessionID uuid.UUID) (model.Session, error)
	UpdateSessionLastUsed(ctx context.Context, sessionID uuid.UUID) error
	GetActiveSessions(ctx context.Context, userID uuid.UUID) ([]model.Session, error)
	GetCurrentSession(ctx context.Context, userID uuid.UUID) (uuid.UUID, error)
	CleanupInactiveSessions(ctx context.Context) error
}
type RoleRepository interface {
	GetAllRoles(ctx context.Context) ([]model.Role, error)
	GetUserRoles(ctx context.Context, userID uuid.UUID) ([]model.Role, error)
	AssignRoleToUser(ctx context.Context, userID uuid.UUID, roleID int64) error
	RemoveRoleFromUser(ctx context.Context, userID uuid.UUID, roleID int64) error
	CreateRole(ctx context.Context, roleName string) (int64, error)
	DeleteRole(ctx context.Context, roleID int64) error
	UpdateRole(ctx context.Context, roleID int64, newRoleName string) error
}
