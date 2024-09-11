package config

import (
	"errors"
	consts "github.com/t34-dev/go-svc-starter/internal"
	"os"
	"time"
)

const (
	refreshTokenSecretKeyEnvName  = "AUTH_REFRESH_TOKEN_SECRET_KEY"
	accessTokenSecretKeyEnvName   = "AUTH_ACCESS_TOKEN_SECRET_KEY"
	refreshTokenExpirationEnvName = "AUTH_REFRESH_TOKEN_EXPIRATION"
	accessTokenExpirationEnvName  = "AUTH_ACCESS_TOKEN_EXPIRATION"
)

var _ AUTHConfig = &authConfig{}

type AUTHConfig interface {
	RefreshTokenSecretKey() string
	AccessTokenSecretKey() string
	RefreshTokenExpiration() time.Duration
	AccessTokenExpiration() time.Duration
}

type authConfig struct {
	refreshTokenSecretKey string
	accessTokenSecretKey  string
}

func (d *authConfig) RefreshTokenSecretKey() string {
	return "a.RefreshTokenSecretKey"
}

func (d *authConfig) AccessTokenSecretKey() string {
	return d.accessTokenSecretKey
}

func (d *authConfig) RefreshTokenExpiration() time.Duration {
	return consts.RefreshTokenExpiration
}

func (d *authConfig) AccessTokenExpiration() time.Duration {
	return consts.AccessTokenExpiration
}

func NewAUTHConfig() (AUTHConfig, error) {
	refreshTokenSecretKey := os.Getenv(refreshTokenSecretKeyEnvName)
	if len(refreshTokenSecretKey) == 0 {
		return nil, errors.New("RefreshTokenSecretKey not found")
	}

	accessTokenSecretKey := os.Getenv(accessTokenSecretKeyEnvName)
	if len(accessTokenSecretKey) == 0 {
		return nil, errors.New("accessTokenSecretKey not found")
	}

	return &authConfig{
		refreshTokenSecretKey: refreshTokenSecretKey,
		accessTokenSecretKey:  accessTokenSecretKey,
	}, nil
}
