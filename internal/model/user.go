package model

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	LogoURL   *string   `json:"logo_url,omitempty"`
	IsBlocked bool      `json:"is_blocked"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Roles     []Role    `json:"roles,omitempty"`
}

type UserInfo struct {
	User     User      `json:"user"`
	Sessions []Session `json:"sessions"`
}
