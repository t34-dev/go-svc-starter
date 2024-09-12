package config

import "net"

var _ GrpcConfig = &grpcConfig{}

type GrpcConfig interface {
	Host() string
	Port() string
	Address() string
}
type grpcConfig struct {
	HostVal string `yaml:"host" env:"GRPC_HOST" yaml-default:"localhost"`
	PortVal string `yaml:"port" env:"GRPC_PORT" yaml-default:"50051"`
}

func (a *grpcConfig) Host() string {
	return a.HostVal
}

func (a *grpcConfig) Port() string {
	return a.PortVal
}
func (a *grpcConfig) Address() string {
	return net.JoinHostPort(a.HostVal, a.PortVal)
}
