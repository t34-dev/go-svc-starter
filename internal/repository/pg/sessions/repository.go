package session_repository

import (
	"context"
	"fmt"
	"github.com/t34-dev/go-svc-starter/internal/model"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/t34-dev/go-svc-starter/internal/repository"
	"github.com/t34-dev/go-utils/pkg/db"
)

const (
	sessionTable       = "sessions"
	idColumn           = "id"
	userIDColumn       = "user_id"
	deviceKeyColumn    = "device_key"
	deviceNameColumn   = "device_name"
	lastUsedColumn     = "last_used"
	refreshTokenColumn = "refresh_token"
	expiresAtColumn    = "expires_at"
	createdAtColumn    = "created_at"
	updatedAtColumn    = "updated_at"
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

func (r sessionRepository) CreateSession(ctx context.Context, userID int64, sessionKey, sessionName, refreshToken string) error {
	query, args, err := r.builder.Insert(sessionTable).
		Columns(userIDColumn, deviceKeyColumn, deviceNameColumn, lastUsedColumn, refreshTokenColumn, expiresAtColumn, createdAtColumn, updatedAtColumn).
		Values(userID, sessionKey, sessionName, time.Now(), refreshToken, time.Now().Add(24*time.Hour), time.Now(), time.Now()).
		ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "session_repository.CreateSession",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	return err
}

func (r sessionRepository) UpsertSession(ctx context.Context, userID int64, sessionKey, sessionName, refreshToken string) error {
	query, args, err := r.builder.Insert(sessionTable).
		Columns(userIDColumn, deviceKeyColumn, deviceNameColumn, lastUsedColumn, refreshTokenColumn, expiresAtColumn, createdAtColumn, updatedAtColumn).
		Values(userID, sessionKey, sessionName, time.Now(), refreshToken, time.Now().Add(24*time.Hour), time.Now(), time.Now()).
		Suffix("ON CONFLICT (user_id, device_key) DO UPDATE SET " +
			deviceNameColumn + " = EXCLUDED." + deviceNameColumn + ", " +
			lastUsedColumn + " = EXCLUDED." + lastUsedColumn + ", " +
			refreshTokenColumn + " = EXCLUDED." + refreshTokenColumn + ", " +
			expiresAtColumn + " = EXCLUDED." + expiresAtColumn + ", " +
			updatedAtColumn + " = EXCLUDED." + updatedAtColumn).
		ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "session_repository.UpsertSession",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	return err
}

func (r sessionRepository) DeleteSession(ctx context.Context, userID int64, refreshToken string) error {
	query, args, err := r.builder.Delete(sessionTable).
		Where(sq.Eq{userIDColumn: userID, refreshTokenColumn: refreshToken}).
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

func (r sessionRepository) GetSessionByRefreshToken(ctx context.Context, refreshToken, sessionKey string) (model.Session, error) {
	query, args, err := r.builder.Select(userIDColumn, expiresAtColumn).
		From(sessionTable).
		Where(sq.Eq{refreshTokenColumn: refreshToken, deviceKeyColumn: sessionKey}).
		Limit(1).
		ToSql()
	if err != nil {
		return model.Session{}, err
	}

	q := db.Query{
		Name:     "session_repository.GetSessionByRefreshToken",
		QueryRaw: query,
	}

	var session model.Session
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&session.UserID, &session.ExpiresAt)
	return session, err
}

func (r sessionRepository) UpdateSession(ctx context.Context, userID int64, sessionKey, sessionName, refreshToken string) error {
	query, args, err := r.builder.Update(sessionTable).
		Set(deviceNameColumn, sessionName).
		Set(lastUsedColumn, time.Now()).
		Set(refreshTokenColumn, refreshToken).
		Set(expiresAtColumn, time.Now().Add(24*time.Hour)).
		Set(updatedAtColumn, time.Now()).
		Where(sq.Eq{userIDColumn: userID, deviceKeyColumn: sessionKey}).
		ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "session_repository.UpdateSession",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	return err
}

func (r sessionRepository) GetActiveSessions(ctx context.Context, userID int64) ([]model.Session, error) {
	query, args, err := r.builder.Select(
		idColumn, userIDColumn, deviceKeyColumn, deviceNameColumn,
		lastUsedColumn, refreshTokenColumn, expiresAtColumn,
		createdAtColumn, updatedAtColumn,
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
			&session.LastUsed, &session.RefreshToken, &session.ExpiresAt,
			&session.CreatedAt, &session.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan session: %v", err)
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}

func (r sessionRepository) GetCurrentSession(ctx context.Context, userID int64) (int64, error) {
	query, args, err := r.builder.Select(idColumn).
		From(sessionTable).
		Where(sq.And{
			sq.Eq{userIDColumn: userID},
			sq.Gt{expiresAtColumn: time.Now()},
		}).
		OrderBy(lastUsedColumn + " DESC").
		Limit(1).
		ToSql()
	if err != nil {
		return 0, err
	}

	q := db.Query{
		Name:     "session_repository.GetCurrentSession",
		QueryRaw: query,
	}

	var id int64
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
