package watch

import (
	"github.com/getnoops/ops/pkg/logging"

	"github.com/spf13/viper"
)

type Config struct {
	DeploymentId string `mapstructure:"deployment"`
	Log          *logging.Config
}

func MustNewConfig(v *viper.Viper) *Config {
	config := new(Config)
	err := v.Unmarshal(config)
	logging.OnError(err).Fatal("unable to read config")

	err = config.Log.SetLogger()
	logging.OnError(err).Fatal("unable to set logger")

	return config
}
