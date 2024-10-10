package service

import (
	"context"
	"github.com/google/uuid"
	othergrpcservice "github.com/t34-dev/go-svc-starter/internal/client/other_grpc_service"
	"github.com/t34-dev/go-svc-starter/internal/model"
	"github.com/t34-dev/go-svc-starter/internal/repository"
	role_manager "github.com/t34-dev/go-svc-starter/pkg/role-manager"
	"github.com/t34-dev/go-utils/pkg/db"
	"time"
)

type Dependencies struct {
	Service       Service
	Repos         repository.Repository
	OtherService  othergrpcservice.OtherGRPCService
	AccessManager role_manager.RoleManager
	TxManager     db.TxManager
}

func NewDeps(service Service,
	repos repository.Repository,
	otherService othergrpcservice.OtherGRPCService,
	accessManager role_manager.RoleManager,
	txManager db.TxManager,
) Dependencies {
	return Dependencies{
		Service:       service,
		Repos:         repos,
		OtherService:  otherService,
		AccessManager: accessManager,
		TxManager:     txManager,
	}
}

type Service struct {
	Common CommonService
	Access AccessService
	Auth   AuthService
}

type CommonService interface {
	GetDBTime(ctx context.Context) (time.Time, error)
	GetTime(ctx context.Context) (time.Time, error)
	GetPost(ctx context.Context, id int64) (*model.Post, error)
}
type AccessService interface {
	Check(ctx context.Context, path string) (bool, error)
}
type AuthService interface {
	Registration(ctx context.Context, email, username, password, deviceKey, deviceName string) (*model.AuthTokens, error)
	Login(ctx context.Context, login, password, deviceKey, deviceName string) (*model.AuthTokens, error)
	Logout(ctx context.Context, token string) error
	GetUserInfo(ctx context.Context, userID uuid.UUID) (*model.UserInfo, error)
	GetActiveSessions(ctx context.Context, userID uuid.UUID) ([]model.Session, error)
	RefreshToken(ctx context.Context, refreshToken string) (*model.AuthTokens, error)
	ValidateToken(ctx context.Context, token string) (*model.ValidateTokenResponse, error)
	RevokeSession(ctx context.Context, sessionID uuid.UUID) error
	GetAllRoles(ctx context.Context) ([]model.Role, error)
	AssignRoleToUser(ctx context.Context, userID uuid.UUID, roleID int64) error
	RemoveRoleFromUser(ctx context.Context, userID uuid.UUID, roleID int64) error
	CreateRole(ctx context.Context, roleName string) (int64, error)
	DeleteRole(ctx context.Context, roleID int64) error
	UpdateRole(ctx context.Context, roleID int64, newRoleName string) error
	CleanupInactiveSessions(ctx context.Context) error
}
