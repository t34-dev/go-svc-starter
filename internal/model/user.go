package model

import "time"

type User struct {
	ID              int64     `json:"id"`
	Email           string    `json:"email"`
	Username        string    `json:"username"`
	Password        string    `json:"-"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	CurrentDeviceID int64     `json:"current_device_id,omitempty"`
}
type UserInfo struct {
	User    User     `json:"user"`
	Devices []Device `json:"devices"`
}
