package config

var _ AppConfig = &appConfig{}

type AppConfig interface {
	Environment() string
	AppName() string
	IsTSL() bool
	IsProduction() bool
	LogLevel() string
}
type appConfig struct {
	EnvironmentVal  string `env:"ENV"`
	AppNameVal      string `env:"APP_NAME"`
	IsTSLVal        bool   `yaml:"is_tsl" env:"IS_TSL" yaml-default:"false"`
	IsProductionVal bool   `yaml:"is_production" env:"IS_PRODUCTION" yaml-default:"true"`
	NameVal         string `yaml:"name" env-default:"ms-sso"`
	LogLevelVal     string `yaml:"log_level" env-default:"info"`
}

func (a *appConfig) Environment() string {
	if a.EnvironmentVal == "" {
		return defaultEnvironment
	}
	return a.EnvironmentVal
}
func (a *appConfig) AppName() string {
	if a.AppNameVal == "" {
		return defaultServiceName
	}
	return a.AppNameVal
}
func (a *appConfig) IsTSL() bool {
	return a.IsTSLVal
}
func (a *appConfig) IsProduction() bool {
	return a.IsProductionVal
}

func (a *appConfig) LogLevel() string {
	return a.LogLevelVal
}
