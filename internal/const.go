package consts

import "time"

const (
	AuthPrefix = "Bearer "

	RefreshTokenExpiration = 60 * time.Minute
	AccessTokenExpiration  = 5 * time.Minute
)
