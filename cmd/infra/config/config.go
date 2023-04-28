package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config struct {
	App string
	Log struct {
		Level string
	}
	Server struct {
		Address string
		Timeout uint
	}

	Database struct {
		Host     string
		Port     string
		User     string
		Name     string
		Password string
		SslMode  string
	}
}

func ProvideConfig(logger *zap.SugaredLogger) Config {
	logger.Info("Executing ProvideConfig.")
	var Config Config
	viper.SetConfigFile(`cloud.config.yaml`)
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logger.Fatal("config file not found", zap.Error(err))
		} else {
			logger.Fatal("error on config", zap.Error(err))
		}

	}
	err = viper.Unmarshal(&Config)
	if err != nil {
		logger.Fatal("error on config", zap.Error(err))
	}
	return Config
}
