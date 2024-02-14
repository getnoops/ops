package info

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/getnoops/ops/pkg/config"
	"github.com/getnoops/ops/pkg/version"
)

type Config struct {
}

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version number of ops",
		Long:  `Will show information about the current version of ops.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			cfg, err := config.New[Config](ctx, viper.GetViper())
			if err != nil {
				return err
			}
			return Info(ctx, cfg)
		},
	}

	return cmd
}

func Info(ctx context.Context, cfg *config.NoOps[Config]) error {
	cfg.WriteStdout("Ops Version " + version.Version() + "\n\n")
	cfg.WriteStdout("Ops API " + cfg.Api.GraphQL)
	return nil
}
