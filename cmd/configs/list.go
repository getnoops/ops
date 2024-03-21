package configs

import (
	"context"

	"github.com/getnoops/ops/pkg/config"
	"github.com/getnoops/ops/pkg/queries"
	"github.com/getnoops/ops/pkg/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ListConfig struct {
	Page     int `mapstructure:"page" default:"1"`
	PageSize int `mapstructure:"page-size" default:"10"`
}

func ListCommand(classes []queries.ConfigClass) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "list",
		Short:  "list projects accessible by the active account",
		PreRun: util.BindPreRun,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			return List(ctx, classes)
		},
	}

	util.BindIntFlag(cmd, "page", "The page to load", 1)
	util.BindIntFlag(cmd, "page-size", "The number of items in the page", 10)
	return cmd
}

func List(ctx context.Context, classes []queries.ConfigClass) error {
	cfg, err := config.New[ListConfig, *queries.ConfigItem](ctx, viper.GetViper())
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

	configs, err := q.GetConfigs(ctx, organisation.Id, classes, cfg.Command.Page, cfg.Command.PageSize)
	if err != nil {
		return err
	}

	cfg.WriteList(configs.Items)
	return nil
}
