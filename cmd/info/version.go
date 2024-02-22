package info

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/getnoops/ops/pkg/config"
	"github.com/getnoops/ops/pkg/util"
	"github.com/getnoops/ops/pkg/version"
)

type Config struct {
}

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "version",
		Short:  "Print the version number of ops",
		Long:   `Will show information about the current version of ops.`,
		PreRun: util.BindPreRun,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			return Info(ctx)
		},
	}

	return cmd
}

func Info(ctx context.Context) error {
	cfg, err := config.New[Config, string](ctx, viper.GetViper())
	if err != nil {
		return err
	}

	cfg.WriteStdout("Ops Version " + version.Version() + "\n\n")
	cfg.WriteStdout("Ops API " + cfg.Api.GraphQL)
	return nil
}
