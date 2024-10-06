package model

import (
	"time"
)

type UserX struct {
	Created        time.Time  `json:"created" db:"created" example:"2021-12-01 15:00:45"`
	Nickname       string     `json:"nickname" binding:"required,min=3,max=50" example:"zak"`
	Email          string     `json:"email" binding:"required,email,max=64" db:"email" example:"root@gmail.com"`
	Roles          RolesX     `json:"roles" db:"roles" example:"root,vip"`
	PasswordHash   string     `json:"-" binding:"required" db:"password_hash" example:"12345678"`
	Id             int        `json:"id" db:"id"`
	EmailConfirmed *time.Time `json:"email_confirmed,omitempty" db:"email_confirmed" example:"2021-12-01T15:00:45Z"`
	Logo           *string    `json:"logo" db:"logo" example:"http://url"`
	IsBlock        bool       `json:"is_block" db:"is_block" example:"false"`
}
