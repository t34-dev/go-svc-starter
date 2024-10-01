package config

import "net"

var _ PrometheusConfig = &prometheusConfig{}

type PrometheusConfig interface {
	Host() string
	Port() string
	Address() string
}
type prometheusConfig struct {
	HostVal string `yaml:"host" env:"PROMETHEUS_HOST" yaml-default:"localhost"`
	PortVal string `yaml:"port" env:"PROMETHEUS_PORT" yaml-default:"2112"`
}

func (a *prometheusConfig) Host() string {
	return a.HostVal
}

func (a *prometheusConfig) Port() string {
	return a.PortVal
}
func (a *prometheusConfig) Address() string {
	return net.JoinHostPort(a.HostVal, a.PortVal)
}
