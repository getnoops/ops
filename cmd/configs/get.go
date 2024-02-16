package configs

import (
	"context"

	"github.com/getnoops/ops/pkg/config"
	"github.com/getnoops/ops/pkg/queries"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type GetConfig struct {
}

func GetCommand(class queries.ConfigClass) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [code]",
		Short: "Get a config",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			code := args[0]

			ctx := cmd.Context()
			return Get(ctx, class, code)
		},
		ValidArgs: []string{"code"},
	}
	return cmd
}

func Get(ctx context.Context, class queries.ConfigClass, code string) error {
	cfg, err := config.New[ListConfig, *queries.Config](ctx, viper.GetViper())
	if err != nil {
		return err
	}

	q, err := queries.New(ctx, cfg)
	if err != nil {
		return err
	}

	organisation, err := q.GetCurrentOrganisation(ctx)
	if err == config.ErrNoOrganisation {
		cfg.WriteStderr("no organisation set")
		return nil
	}
	if err != nil {
		return err
	}

	config, err := q.GetConfig(ctx, organisation.Id, code)
	if err != nil {
		cfg.WriteStderr("failed to get configs")
		return nil
	}

	if config == nil || config.Class != class {
		cfg.WriteStderr("config not found")
		return nil
	}

	cfg.WriteObject(config)
	return nil
}
