package config

var _ JwtConfig = &jwtConfig{}

type JwtConfig interface {
	Secret() string
}
type jwtConfig struct {
	SecretVal string `env:"JWT_SECRET_KEY" env-required:"true"`
}

func (a *jwtConfig) Secret() string {
	return a.SecretVal
}
