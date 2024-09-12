package config

import "fmt"

var _ PgConfig = &pgConfig{}

type PgConfig interface {
	Host() string
	Port() string
	DBName() string
	User() string
	Password() string
	SSLMode() string
	MinConns() int32
	MaxConns() int32
	DSN() string
}
type pgConfig struct {
	HostVal     string `env:"PG_HOST" env-required:"true"`
	PortVal     string `env:"PG_PORT" env-required:"5432"`
	DBNameVal   string `env:"PG_NAME" env-required:"true"`
	UserVal     string `env:"PG_USER" env-required:"true"`
	PasswordVal string `env:"PG_PASSWORD" env-required:"true"`
	SSLModeVal  string `yaml:"ssl_mode" env:"PG_SSL_MODE" env-default:"disable"`
	MinConnsVal int32  `yaml:"min_conn" env:"PG_MIN_CONN" env-default:"5"`
	MaxConnsVal int32  `yaml:"max_conn" env:"PG_MAX_CONN" env-default:"10"`
}

func (p pgConfig) Host() string {
	return p.HostVal
}

func (p pgConfig) Port() string {
	return p.PortVal
}

func (p pgConfig) DBName() string {
	return p.DBNameVal
}

func (p pgConfig) User() string {
	return p.UserVal
}

func (p pgConfig) Password() string {
	return p.PasswordVal
}

func (p pgConfig) SSLMode() string {
	return p.SSLModeVal
}

func (p pgConfig) MinConns() int32 {
	return p.MinConnsVal
}

func (p pgConfig) MaxConns() int32 {
	return p.MaxConnsVal
}

func (p *pgConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%v",
		p.HostVal, p.PortVal, p.DBNameVal, p.UserVal, p.PasswordVal, p.SSLModeVal)
}
