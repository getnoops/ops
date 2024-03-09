package cmd

import (
	"bytes"
	_ "embed"
	"io"
	"log"
	"strings"

	"github.com/getnoops/ops/cmd/configs"
	"github.com/getnoops/ops/cmd/containerrepository"
	"github.com/getnoops/ops/cmd/deploy"
	"github.com/getnoops/ops/cmd/envs"
	"github.com/getnoops/ops/cmd/info"
	"github.com/getnoops/ops/cmd/keys"
	"github.com/getnoops/ops/cmd/login"
	"github.com/getnoops/ops/cmd/orgs"
	"github.com/getnoops/ops/cmd/secrets"
	"github.com/getnoops/ops/cmd/settings"
	"github.com/getnoops/ops/cmd/this"
	"github.com/getnoops/ops/cmd/upgrade"
	"github.com/getnoops/ops/pkg/queries"
	"github.com/getnoops/ops/pkg/util"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	//go:embed defaults.yaml
	defaultConfig []byte

	configFiles []string
)

func New(out io.Writer, in io.Reader, args []string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ops",
		Short: "The No_Ops cli used to manage deployments",
		PreRun: func(cmd *cobra.Command, args []string) {
			cmd.Flags().VisitAll(func(flag *pflag.Flag) {
				viper.BindPFlag("command."+flag.Name, flag)
			})
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("NOOPS")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetConfigType("yaml")

	if err := viper.ReadConfig(bytes.NewBuffer(defaultConfig)); err != nil {
		log.Fatal(err)
	}

	util.BindStringPersistentFlag(cmd, "organisation", "The organisation to use", "")
	util.BindStringPersistentFlag(cmd, "token", "The token to use", "")
	util.BindStringPersistentFlag(cmd, "format", "The format for printing output", "table")

	cmd.AddCommand(
		info.New(),
		login.New(),
		upgrade.New(),
		settings.New(),
		orgs.New(),
		envs.New(),
		configs.New("Compute", queries.ConfigClassCompute),
		configs.New("Storage", queries.ConfigClassStorage),
		configs.New("Integration", queries.ConfigClassNotification, queries.ConfigClassQueue),
		containerrepository.New(),
		secrets.New(),
		keys.New(),
		deploy.New(),
		this.New(),
	)
	cmd.InitDefaultVersionFlag()
	return cmd
}
