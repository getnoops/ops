package deploy

import (
	"github.com/getnoops/ops/pkg/logging"

	"github.com/spf13/viper"
)

type Config struct {
	Environment string `mapstructure:"env"`
	StackFile   string `mapstructure:"file"`
	Log         *logging.Config
}

func MustNewConfig(v *viper.Viper) *Config {
	config := new(Config)
	err := v.Unmarshal(config)
	logging.OnError(err).Fatal("unable to read config")

	err = config.Log.SetLogger()
	logging.OnError(err).Fatal("unable to set logger")

	return config
}
