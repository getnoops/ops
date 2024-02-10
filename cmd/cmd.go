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
	"github.com/getnoops/ops/cmd/keys"
	"github.com/getnoops/ops/cmd/login"
	"github.com/getnoops/ops/cmd/orgs"
	"github.com/getnoops/ops/cmd/settings"
	"github.com/getnoops/ops/cmd/upgrade"
	"github.com/getnoops/ops/pkg/queries"
	"github.com/getnoops/ops/pkg/util"
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

	util.BindStringPersistentFlag(cmd, "organisation", "The organisation to use", "")
	util.BindStringPersistentFlag(cmd, "format", "The format for printing output", "table")

	cmd.AddCommand(
		login.New(),
		upgrade.New(),
		settings.New(),
		orgs.New(),
		envs.New(),
		configs.New("Compute", queries.ConfigClassCompute),
		configs.New("Storage", queries.ConfigClassStorage),
		configs.New("Integration", queries.ConfigClassIntegration),
		containerrepository.New(),
		keys.New(),
		deploy.New(),
	)
	cmd.InitDefaultVersionFlag()
	return cmd
}
