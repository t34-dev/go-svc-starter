package model

import "time"

type User struct {
	ID               int64     `json:"id"`
	Email            string    `json:"email"`
	Username         string    `json:"username"`
	Password         string    `json:"-"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	CurrentSessionID int64     `json:"current_session_id,omitempty"`
}
type UserInfo struct {
	User     User      `json:"user"`
	Sessions []Session `json:"sessions"`
}
