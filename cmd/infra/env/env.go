package env

import (
	"github.com/spf13/viper"
)

type config struct {
	Server struct {
		Address string
		Timeout int64
	}
}

var Config config

func init() {
	viper.SetConfigFile(`config.yaml`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&Config)
	if err != nil {
		panic(err)
	}
}
