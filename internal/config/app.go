package config

import (
	"go.uber.org/zap/zapcore"
)

var _ AppConfig = &appConfig{}

type AppConfig interface {
	IsProduction() bool
	Name() string
	LogLevel() zapcore.Level
}
type appConfig struct {
	IsProductionVal bool          `yaml:"is_production" env:"IS_PRODUCTION" yaml-default:"true"`
	NameVal         string        `yaml:"name" env-default:"ms-sso"`
	LogLevelVal     zapcore.Level `yaml:"log_level" env-default:"info"`
}

func (a *appConfig) IsProduction() bool {
	return a.IsProductionVal
}

func (a *appConfig) Name() string {
	return a.NameVal
}

func (a *appConfig) LogLevel() zapcore.Level {
	return a.LogLevelVal
}
