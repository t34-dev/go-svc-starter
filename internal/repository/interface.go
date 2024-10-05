package repository

import (
	"context"
	"github.com/t34-dev/go-svc-starter/internal/model"
	"time"
)

type Repository struct {
	Common CommonRepository
	User   UserRepository
	Device DeviceRepository
}

type CommonRepository interface {
	GetTime(ctx context.Context) (time.Time, error)
}
type UserRepository interface {
	CreateUser(ctx context.Context, email, username, hashedPassword string) (int64, error)
	GetUserByLogin(ctx context.Context, login string) (model.User, error)
	GetUserInfo(ctx context.Context, userID int64) (model.User, error)
}
type DeviceRepository interface {
	CreateDevice(ctx context.Context, userID int64, deviceKey, deviceName, refreshToken string) error
	UpsertDevice(ctx context.Context, userID int64, deviceKey, deviceName, refreshToken string) error
	DeleteDevice(ctx context.Context, userID int64, refreshToken string) error
	GetDeviceByRefreshToken(ctx context.Context, refreshToken, deviceKey string) (model.Device, error)
	UpdateDevice(ctx context.Context, userID int64, deviceKey, deviceName, refreshToken string) error
	GetActiveDevices(ctx context.Context, userID int64) ([]model.Device, error)
	GetCurrentDevice(ctx context.Context, userID int64) (int64, error)
	CleanupInactiveSessions(ctx context.Context) error
}
