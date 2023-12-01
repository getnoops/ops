package cmd

import (
	"bytes"
	_ "embed"
	"io"
	"log"
	"strings"

	"github.com/getnoops/ops/cmd/login"
	"github.com/getnoops/ops/cmd/upgrade"
	"github.com/getnoops/ops/pkg/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	//go:embed defaults.yaml
	defaultConfig []byte

	configFiles []string
)

func New(out io.Writer, in io.Reader, args []string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ops",
		Short:   "The No_Ops cli used to manage deployments",
		Version: version.Version(),
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("NOOPS")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetConfigType("yaml")
	if err := viper.ReadConfig(bytes.NewBuffer(defaultConfig)); err != nil {
		log.Fatal(err)
	}

	cobra.OnInitialize(initConfig)
	cmd.PersistentFlags().StringArrayVar(&configFiles, "config", nil, "path to config file to overwrite system defaults")

	cmd.AddCommand(
		login.New(),
		upgrade.New(),
	)

	cmd.InitDefaultVersionFlag()
	return cmd
}

func initConfig() {
	for _, file := range configFiles {
		viper.SetConfigFile(file)
		if err := viper.MergeInConfig(); err != nil {
			log.Fatal(err)
		}
	}
}
