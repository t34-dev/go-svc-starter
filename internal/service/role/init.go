package role_service

import (
	"context"
	"github.com/google/uuid"
	"github.com/t34-dev/go-svc-starter/internal/model"
	"github.com/t34-dev/go-svc-starter/internal/service"
)

var _ service.RoleService = &roleService{}

type roleService struct {
	opt service.Options
}

func New(opt service.Options) service.RoleService {
	return &roleService{
		opt: opt,
	}
}

func (s *roleService) GetAllRoles(ctx context.Context) ([]model.Role, error) {
	return s.opt.Repos.Role.GetAllRoles(ctx)
}

func (s *roleService) AddRoleToUser(ctx context.Context, userID uuid.UUID, roleID int64) error {
	return s.opt.Repos.Role.AddRoleToUser(ctx, userID, roleID)
}

func (s *roleService) RemoveRoleFromUser(ctx context.Context, userID uuid.UUID, roleID int64) error {
	return s.opt.Repos.Role.RemoveRoleFromUser(ctx, userID, roleID)
}

func (s *roleService) CreateRole(ctx context.Context, roleName string) (int64, error) {
	return s.opt.Repos.Role.CreateRole(ctx, roleName)
}

func (s *roleService) DeleteRole(ctx context.Context, roleID int64) error {
	return s.opt.Repos.Role.DeleteRole(ctx, roleID)
}

func (s *roleService) UpdateRole(ctx context.Context, roleID int64, newRoleName string) error {
	return s.opt.Repos.Role.UpdateRole(ctx, roleID, newRoleName)
}
