package auth

import (
	"github.com/getnoops/ops/pkg/logging"

	"github.com/spf13/viper"
)

type AuthConfig struct {
	Issuer   string
	ClientId string
	Scopes   []string
}

type Config struct {
	Company string `mapstructure:"company"`
	Auth    AuthConfig
	Log     *logging.Config
}

func MustNewConfig(v *viper.Viper) *Config {
	config := new(Config)
	err := v.Unmarshal(config)
	logging.OnError(err).Fatal("unable to read config")

	err = config.Log.SetLogger()
	logging.OnError(err).Fatal("unable to set logger")

	return config
}
