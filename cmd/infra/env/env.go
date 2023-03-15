package env

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type config struct {
	Server struct {
		Address string
		Timeout int64
	}

	Database struct {
		Host     string
		Port     string
		User     string
		Name     string
		Password string
	}
}

var Config config

func init() {
	viper.SetConfigFile(`cloud.config.yaml`)
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			zap.L().Fatal("config file not found", zap.Error(err))
		} else {
			zap.L().Fatal("error on config", zap.Error(err))
		}

	}
	err = viper.Unmarshal(&Config)
	if err != nil {
		zap.L().Fatal("error on config", zap.Error(err))
	}
}
