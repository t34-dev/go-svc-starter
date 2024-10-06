package model

import "time"

type Session struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	DeviceKey    string    `json:"device_key"`
	DeviceName   string    `json:"device_name"`
	LastUsed     time.Time `json:"last_used"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
