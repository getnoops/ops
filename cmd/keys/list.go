package keys

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

func ListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "list",
		Short:  "list api keys",
		PreRun: util.BindPreRun,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			return List(ctx)
		},
	}

	util.BindIntFlag(cmd, "page", "The page to load", 1)
	util.BindIntFlag(cmd, "page-size", "The number of items in the page", 10)
	return cmd
}

func List(ctx context.Context) error {
	cfg, err := config.New[ListConfig, *queries.ApiKey](ctx, viper.GetViper())
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

	out, err := q.GetApiKeys(ctx, organisation.Id, cfg.Command.Page, cfg.Command.PageSize)
	if err != nil {
		cfg.WriteStderr("failed to get config")
		return err
	}

	cfg.WriteList(out.Items)
	return nil
}
