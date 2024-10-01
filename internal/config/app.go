package config

var _ AppConfig = &appConfig{}

type AppConfig interface {
	IsTSL() bool
	IsProduction() bool
	Name() string
	LogLevel() string
}
type appConfig struct {
	IsTSLVal        bool   `yaml:"is_tsl" env:"IS_TSL" yaml-default:"false"`
	IsProductionVal bool   `yaml:"is_production" env:"IS_PRODUCTION" yaml-default:"true"`
	NameVal         string `yaml:"name" env-default:"ms-sso"`
	LogLevelVal     string `yaml:"log_level" env-default:"info"`
}

func (a *appConfig) IsTSL() bool {
	return a.IsTSLVal
}
func (a *appConfig) IsProduction() bool {
	return a.IsProductionVal
}

func (a *appConfig) Name() string {
	return a.NameVal
}

func (a *appConfig) LogLevel() string {
	return a.LogLevelVal
}
