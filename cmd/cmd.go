package cmd

import (
	"bytes"
	_ "embed"
	"io"
	"strings"

	"github.com/getnoops/ops/cmd/auth"
	"github.com/getnoops/ops/cmd/deploy"
	"github.com/getnoops/ops/cmd/upgrade"
	"github.com/getnoops/ops/pkg/logging"
	"github.com/getnoops/ops/pkg/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//go:embed defaults.yaml
var defaultConfig []byte

var configFiles []string

func New(out io.Writer, in io.Reader, args []string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ops",
		Short:   "The No_Ops cli used to manage deployments",
		Version: version.Version(),
	}

	// This is just temporary. We should read this value from .env later when deploying CLI/Brain to other environments.
	viper.SetDefault("BrainUrl", "http://localhost:8080")

	viper.AutomaticEnv()
	viper.SetEnvPrefix("NOOPS")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetConfigType("yaml")
	err := viper.ReadConfig(bytes.NewBuffer(defaultConfig))
	logging.OnError(err).Fatal("unable to read default config")

	cobra.OnInitialize(initConfig)
	cmd.PersistentFlags().StringArrayVar(&configFiles, "config", nil, "path to config file to overwrite system defaults")

	cmd.AddCommand(
		auth.New(),
		upgrade.New(),
		deploy.New(),
	)

	cmd.InitDefaultVersionFlag()
	return cmd
}

func initConfig() {
	for _, file := range configFiles {
		viper.SetConfigFile(file)
		err := viper.MergeInConfig()
		logging.WithFields("file", file).OnError(err).Warn("unable to read config file")
	}
}
