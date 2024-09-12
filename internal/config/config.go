package config

import (
	"context"
	"github.com/google/go-cmp/cmp"
	"github.com/t34-dev/go-svc-starter/pkg/iconfig"
)

var cfg *config

type result struct {
	newConfig *config
	Error     error
}

var resultChan = make(chan result, 1)

type config struct {
	App  appConfig  `yaml:"app"`
	Grpc grpcConfig `yaml:"grpc"`
	Http httpConfig `yaml:"http"`
	Pg   pgConfig   `yaml:"pg"`
}

func App() AppConfig {
	return &cfg.App
}

func Grpc() GrpcConfig {
	return &cfg.Grpc
}

func Http() HttpConfig {
	return &cfg.Http
}
func Pg() PgConfig {
	return &cfg.Pg
}

func New(ctx context.Context, yamlPath, envPath string) (error, <-chan result) {
	if cfg != nil {
		return nil, resultChan
	}
	if nil == cfg {
		defer func() {
			err := watchConfig(ctx, yamlPath, envPath, func(newConfig *config, err error) {
				resultChan <- result{newConfig: newConfig, Error: err}
			})
			if err != nil {
				resultChan <- result{newConfig: nil, Error: err}
			}
		}()
	}

	initialConfig := new(config)
	err := iconfig.GetConfig(initialConfig, yamlPath, envPath)
	if err != nil {
		return err, resultChan
	}
	cfg = initialConfig

	return nil, resultChan
}

func watchConfig(ctx context.Context, yamlPath string, envPath string, callBack func(newConfig *config, err error)) error {
	return iconfig.WatchConfig(ctx, yamlPath, func(msg string, err error) {
		if err != nil {
			callBack(nil, err)
		} else {
			oldConfig := cfg

			newConfig := new(config)
			err := iconfig.GetConfig(newConfig, yamlPath, envPath)
			if err != nil {
				callBack(nil, err)
			} else {
				if isChangedConfig(oldConfig, newConfig) {
					cfg = newConfig
					callBack(newConfig, nil)
				}
			}
		}
	})
}
func isChangedConfig(oldConfig, newConfig *config) bool {
	if oldConfig == nil || newConfig == nil {
		return true
	}

	diff := cmp.Diff(oldConfig, newConfig)
	return diff != ""
}
