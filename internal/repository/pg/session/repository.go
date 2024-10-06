package session_repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/t34-dev/go-svc-starter/internal/model"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/t34-dev/go-svc-starter/internal/repository"
	"github.com/t34-dev/go-utils/pkg/db"
)

const (
	sessionTable     = "sessions"
	idColumn         = "id"
	userIDColumn     = "user_id"
	deviceKeyColumn  = "device_key"
	deviceNameColumn = "device_name"
	expiresAtColumn  = "expires_at"
	createdAtColumn  = "created_at"
	updatedAtColumn  = "updated_at"
)

var _ repository.SessionRepository = (*sessionRepository)(nil)

type sessionRepository struct {
	db      db.Client
	builder sq.StatementBuilderType
}

func New(db db.Client) repository.SessionRepository {
	return &sessionRepository{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r sessionRepository) CreateSession(ctx context.Context, userID uuid.UUID, deviceKey, deviceName string) (uuid.UUID, error) {
	sessionID := uuid.New()
	query, args, err := r.builder.Insert(sessionTable).
		Columns(idColumn, userIDColumn, deviceKeyColumn, deviceNameColumn, expiresAtColumn, createdAtColumn, updatedAtColumn).
		Values(sessionID, userID, deviceKey, deviceName, time.Now().Add(24*time.Hour), time.Now(), time.Now()).
		ToSql()
	if err != nil {
		return uuid.Nil, err
	}

	q := db.Query{
		Name:     "session_repository.CreateSession",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return uuid.Nil, err
	}
	return sessionID, nil
}

func (r sessionRepository) UpsertSession(ctx context.Context, userID uuid.UUID, deviceKey, deviceName string) (uuid.UUID, error) {
	sessionID := uuid.New()
	query, args, err := r.builder.Insert(sessionTable).
		Columns(idColumn, userIDColumn, deviceKeyColumn, deviceNameColumn, expiresAtColumn, createdAtColumn, updatedAtColumn).
		Values(sessionID, userID, deviceKey, deviceName, time.Now().Add(24*time.Hour), time.Now(), time.Now()).
		Suffix("ON CONFLICT (user_id, device_key) DO UPDATE SET " +
			deviceNameColumn + " = EXCLUDED." + deviceNameColumn + ", " +
			expiresAtColumn + " = EXCLUDED." + expiresAtColumn + ", " +
			updatedAtColumn + " = EXCLUDED." + updatedAtColumn + " RETURNING " + idColumn).
		ToSql()
	if err != nil {
		return uuid.Nil, err
	}

	q := db.Query{
		Name:     "session_repository.UpsertSession",
		QueryRaw: query,
	}

	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&sessionID)
	if err != nil {
		return uuid.Nil, err
	}
	return sessionID, nil
}

func (r sessionRepository) DeleteSession(ctx context.Context, sessionID uuid.UUID) error {
	query, args, err := r.builder.Delete(sessionTable).
		Where(sq.Eq{idColumn: sessionID}).
		ToSql()
	if err != nil {
		return err
	}
	q := db.Query{
		Name:     "session_repository.DeleteSession",
		QueryRaw: query,
	}
	_, err = r.db.DB().ExecContext(ctx, q, args...)
	return err
}

func (r sessionRepository) GetSessionByID(ctx context.Context, sessionID uuid.UUID) (model.Session, error) {
	query, args, err := r.builder.Select(
		idColumn, userIDColumn, deviceKeyColumn, deviceNameColumn,
		expiresAtColumn, createdAtColumn, updatedAtColumn,
	).From(sessionTable).
		Where(sq.Eq{idColumn: sessionID}).
		Limit(1).
		ToSql()
	if err != nil {
		return model.Session{}, err
	}

	q := db.Query{
		Name:     "session_repository.GetSessionByID",
		QueryRaw: query,
	}

	var session model.Session
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(
		&session.ID, &session.UserID, &session.DeviceKey, &session.DeviceName,
		&session.ExpiresAt, &session.CreatedAt, &session.UpdatedAt,
	)
	return session, err
}

func (r sessionRepository) UpdateSessionLastUsed(ctx context.Context, sessionID uuid.UUID) error {
	query, args, err := r.builder.Update(sessionTable).
		Set(updatedAtColumn, time.Now()).
		Where(sq.Eq{idColumn: sessionID}).
		ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "session_repository.UpdateSessionLastUsed",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	return err
}

func (r sessionRepository) GetActiveSessions(ctx context.Context, userID uuid.UUID) ([]model.Session, error) {
	query, args, err := r.builder.Select(
		idColumn, userIDColumn, deviceKeyColumn, deviceNameColumn,
		expiresAtColumn, createdAtColumn, updatedAtColumn,
	).From(sessionTable).
		Where(sq.And{
			sq.Eq{userIDColumn: userID},
			sq.Gt{expiresAtColumn: time.Now()},
		}).
		ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "session_repository.GetActiveSessions",
		QueryRaw: query,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []model.Session
	for rows.Next() {
		var session model.Session
		err := rows.Scan(
			&session.ID, &session.UserID, &session.DeviceKey, &session.DeviceName,
			&session.ExpiresAt, &session.CreatedAt, &session.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan session: %v", err)
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}

func (r sessionRepository) GetCurrentSession(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	query, args, err := r.builder.Select(idColumn).
		From(sessionTable).
		Where(sq.And{
			sq.Eq{userIDColumn: userID},
			sq.Gt{expiresAtColumn: time.Now()},
		}).
		OrderBy(updatedAtColumn + " DESC").
		Limit(1).
		ToSql()
	if err != nil {
		return uuid.Nil, err
	}

	q := db.Query{
		Name:     "session_repository.GetCurrentSession",
		QueryRaw: query,
	}

	var id uuid.UUID
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&id)
	return id, err
}

func (r sessionRepository) CleanupInactiveSessions(ctx context.Context) error {
	query, args, err := r.builder.Delete(sessionTable).
		Where(sq.Lt{expiresAtColumn: time.Now()}).
		ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "session_repository.CleanupInactiveSessions",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	return err
}
