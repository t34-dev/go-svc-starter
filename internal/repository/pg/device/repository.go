package device_repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/t34-dev/go-svc-starter/internal/model"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/t34-dev/go-svc-starter/internal/repository"
)

const (
	deviceTable              = "devices"
	deviceIDColumn           = "id"
	deviceUserIDColumn       = "user_id"
	deviceKeyColumn          = "device_key"
	deviceNameColumn         = "device_name"
	deviceLastUsedColumn     = "last_used"
	deviceRefreshTokenColumn = "refresh_token"
	deviceExpiresAtColumn    = "expires_at"
	deviceCreatedAtColumn    = "created_at"
	deviceUpdatedAtColumn    = "updated_at"
)

var _ repository.DeviceRepository = (*deviceRepository)(nil)

type deviceRepository struct {
	db      *sql.DB
	builder sq.StatementBuilderType
}

func New(db *sql.DB) repository.DeviceRepository {
	return &deviceRepository{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r deviceRepository) CreateDevice(ctx context.Context, userID int64, deviceKey, deviceName, refreshToken string) error {
	_, err := r.builder.Insert(deviceTable).
		Columns(deviceUserIDColumn, deviceKeyColumn, deviceNameColumn, deviceLastUsedColumn, deviceRefreshTokenColumn, deviceExpiresAtColumn, deviceCreatedAtColumn, deviceUpdatedAtColumn).
		Values(userID, deviceKey, deviceName, time.Now(), refreshToken, time.Now().Add(24*time.Hour), time.Now(), time.Now()).
		RunWith(r.db).
		Exec()
	return err
}

func (r deviceRepository) UpsertDevice(ctx context.Context, userID int64, deviceKey, deviceName, refreshToken string) error {
	_, err := r.builder.Insert(deviceTable).
		Columns(deviceUserIDColumn, deviceKeyColumn, deviceNameColumn, deviceLastUsedColumn, deviceRefreshTokenColumn, deviceExpiresAtColumn, deviceCreatedAtColumn, deviceUpdatedAtColumn).
		Values(userID, deviceKey, deviceName, time.Now(), refreshToken, time.Now().Add(24*time.Hour), time.Now(), time.Now()).
		Suffix("ON CONFLICT (user_id, device_key) DO UPDATE SET " +
			deviceNameColumn + " = EXCLUDED." + deviceNameColumn + ", " +
			deviceLastUsedColumn + " = EXCLUDED." + deviceLastUsedColumn + ", " +
			deviceRefreshTokenColumn + " = EXCLUDED." + deviceRefreshTokenColumn + ", " +
			deviceExpiresAtColumn + " = EXCLUDED." + deviceExpiresAtColumn + ", " +
			deviceUpdatedAtColumn + " = EXCLUDED." + deviceUpdatedAtColumn).
		RunWith(r.db).
		Exec()
	return err
}

func (r deviceRepository) DeleteDevice(ctx context.Context, userID int64, refreshToken string) error {
	_, err := r.builder.Delete(deviceTable).
		Where(sq.Eq{deviceUserIDColumn: userID, deviceRefreshTokenColumn: refreshToken}).
		RunWith(r.db).
		Exec()
	return err
}

func (r deviceRepository) GetDeviceByRefreshToken(ctx context.Context, refreshToken, deviceKey string) (model.Device, error) {
	var device model.Device
	err := r.builder.Select(deviceUserIDColumn, deviceExpiresAtColumn).
		From(deviceTable).
		Where(sq.Eq{deviceRefreshTokenColumn: refreshToken, deviceKeyColumn: deviceKey}).
		Limit(1).
		RunWith(r.db).
		QueryRow().
		Scan(&device.UserID, &device.ExpiresAt)
	return device, err
}

func (r deviceRepository) UpdateDevice(ctx context.Context, userID int64, deviceKey, deviceName, refreshToken string) error {
	_, err := r.builder.Update(deviceTable).
		Set(deviceNameColumn, deviceName).
		Set(deviceLastUsedColumn, time.Now()).
		Set(deviceRefreshTokenColumn, refreshToken).
		Set(deviceExpiresAtColumn, time.Now().Add(24*time.Hour)).
		Set(deviceUpdatedAtColumn, time.Now()).
		Where(sq.Eq{deviceUserIDColumn: userID, deviceKeyColumn: deviceKey}).
		RunWith(r.db).
		Exec()
	return err
}

func (r deviceRepository) GetActiveDevices(ctx context.Context, userID int64) ([]model.Device, error) {
	rows, err := r.builder.Select(
		deviceIDColumn, deviceUserIDColumn, deviceKeyColumn, deviceNameColumn,
		deviceLastUsedColumn, deviceRefreshTokenColumn, deviceExpiresAtColumn,
		deviceCreatedAtColumn, deviceUpdatedAtColumn,
	).From(deviceTable).
		Where(sq.And{
			sq.Eq{deviceUserIDColumn: userID},
			sq.Gt{deviceExpiresAtColumn: time.Now()},
		}).
		RunWith(r.db).
		Query()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []model.Device
	for rows.Next() {
		var device model.Device
		err := rows.Scan(
			&device.ID, &device.UserID, &device.DeviceKey, &device.DeviceName,
			&device.LastUsed, &device.RefreshToken, &device.ExpiresAt,
			&device.CreatedAt, &device.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan device: %v", err)
		}
		devices = append(devices, device)
	}

	return devices, nil
}

func (r deviceRepository) GetCurrentDevice(ctx context.Context, userID int64) (int64, error) {
	var id int64
	err := r.builder.Select(deviceIDColumn).
		From(deviceTable).
		Where(sq.And{
			sq.Eq{deviceUserIDColumn: userID},
			sq.Gt{deviceExpiresAtColumn: time.Now()},
		}).
		OrderBy(deviceLastUsedColumn + " DESC").
		Limit(1).
		RunWith(r.db).
		QueryRow().
		Scan(&id)
	return id, err
}

func (r deviceRepository) CleanupInactiveSessions(ctx context.Context) error {
	_, err := r.builder.Delete(deviceTable).
		Where(sq.Lt{deviceExpiresAtColumn: time.Now()}).
		RunWith(r.db).
		Exec()
	return err
}
