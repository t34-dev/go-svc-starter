package config

import "net"

var _ HttpConfig = &httpConfig{}

type HttpConfig interface {
	Host() string
	Port() string
	Address() string
}
type httpConfig struct {
	HostVal string `yaml:"host" env:"HTTP_HOST" yaml-default:"localhost"`
	PortVal string `yaml:"port" env:"HTTP_PORT" yaml-default:"8080"`
}

func (a *httpConfig) Host() string {
	return a.HostVal
}

func (a *httpConfig) Port() string {
	return a.PortVal
}
func (a *httpConfig) Address() string {
	return net.JoinHostPort(a.HostVal, a.PortVal)
}
