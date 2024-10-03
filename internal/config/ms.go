package config

var _ MSConfig = &msConfig{}

type MSConfig interface {
	OtherGrpcAddress() string
}
type msConfig struct {
	OtherGrpcAddressVal string `yaml:"other_grpc_address" env:"MS_OTHER_GRPC_ADDRESS" env-required:"true"`
}

func (a *msConfig) OtherGrpcAddress() string {
	return a.OtherGrpcAddressVal
}
