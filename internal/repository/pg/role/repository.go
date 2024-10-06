package role_repository

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/t34-dev/go-svc-starter/internal/model"
	"github.com/t34-dev/go-svc-starter/internal/repository"
	"github.com/t34-dev/go-utils/pkg/db"
)

const (
	roleTable     = "roles"
	userRoleTable = "user_roles"
	idColumn      = "id"
	roleColumn    = "role"
	userIDColumn  = "user_id"
	roleIDColumn  = "role_id"
)

var _ repository.RoleRepository = (*roleRepository)(nil)

type roleRepository struct {
	db      db.Client
	builder sq.StatementBuilderType
}

func New(db db.Client) repository.RoleRepository {
	return &roleRepository{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}
func (r roleRepository) GetAllRoles(ctx context.Context) ([]model.Role, error) {
	query, args, err := r.builder.Select(idColumn, roleColumn).
		From(roleTable).
		ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "role_repository.GetAllRoles",
		QueryRaw: query,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []model.Role
	for rows.Next() {
		var role model.Role
		if err := rows.Scan(&role.ID, &role.Name); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, nil
}

func (r roleRepository) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]model.Role, error) {
	query, args, err := r.builder.Select("r."+idColumn, "r."+roleColumn).
		From(roleTable + " r").
		Join(userRoleTable + " ur ON r.id = ur.role_id").
		Where(sq.Eq{"ur.user_id": userID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "role_repository.GetUserRoles",
		QueryRaw: query,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []model.Role
	for rows.Next() {
		var role model.Role
		if err := rows.Scan(&role.ID, &role.Name); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, nil
}

func (r roleRepository) AssignRoleToUser(ctx context.Context, userID uuid.UUID, roleID int64) error {
	query, args, err := r.builder.Insert(userRoleTable).
		Columns(userIDColumn, roleIDColumn).
		Values(userID, roleID).
		Suffix("ON CONFLICT DO NOTHING").
		ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "role_repository.AssignRoleToUser",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	return err
}

func (r roleRepository) RemoveRoleFromUser(ctx context.Context, userID uuid.UUID, roleID int64) error {
	query, args, err := r.builder.Delete(userRoleTable).
		Where(sq.Eq{userIDColumn: userID, roleIDColumn: roleID}).
		ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "role_repository.RemoveRoleFromUser",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	return err
}

func (r roleRepository) CreateRole(ctx context.Context, roleName string) (int64, error) {
	query, args, err := r.builder.Insert(roleTable).
		Columns(roleColumn).
		Values(roleName).
		Suffix("RETURNING " + idColumn).
		ToSql()
	if err != nil {
		return 0, err
	}

	q := db.Query{
		Name:     "role_repository.CreateRole",
		QueryRaw: query,
	}

	var id int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&id)
	return id, err
}

func (r roleRepository) DeleteRole(ctx context.Context, roleID int64) error {
	query, args, err := r.builder.Delete(roleTable).
		Where(sq.Eq{idColumn: roleID}).
		ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "role_repository.DeleteRole",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	return err
}

func (r roleRepository) UpdateRole(ctx context.Context, roleID int64, newRoleName string) error {
	query, args, err := r.builder.Update(roleTable).
		Set(roleColumn, newRoleName).
		Where(sq.Eq{idColumn: roleID}).
		ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "role_repository.UpdateRole",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	return err
}
