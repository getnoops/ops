package configs

import (
	"context"
	"fmt"

	"github.com/getnoops/ops/pkg/config"
	"github.com/getnoops/ops/pkg/queries"
	"github.com/getnoops/ops/pkg/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type GetConfig struct {
}

func GetCommand(classes []queries.ConfigClass) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "get [code]",
		Short:  "Get a config",
		Args:   cobra.ExactArgs(1),
		PreRun: util.BindPreRun,
		RunE: func(cmd *cobra.Command, args []string) error {
			code := args[0]

			ctx := cmd.Context()
			return Get(ctx, classes, code)
		},
		ValidArgs: []string{"code"},
	}
	return cmd
}

func Get(ctx context.Context, classes []queries.ConfigClass, code string) error {
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
		return err
	}

	if config == nil {
		cfg.WriteStderr(fmt.Sprintf("config '%v' was not found", code))
		return nil
	}

	cfg.WriteObject(config)
	return nil
}
