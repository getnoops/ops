package containerrepository

import (
	"context"

	"github.com/getnoops/ops/pkg/config"
	"github.com/getnoops/ops/pkg/queries"
	"github.com/getnoops/ops/pkg/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ListConfig struct {
}

func ListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "list [compute]",
		Short:  "list container registries",
		Args:   cobra.ExactArgs(1),
		PreRun: util.BindPreRun,
		RunE: func(cmd *cobra.Command, args []string) error {
			configCode := args[0]

			ctx := cmd.Context()
			return List(ctx, configCode)
		},
	}
	return cmd
}

func List(ctx context.Context, configCode string) error {
	cfg, err := config.New[ListConfig, *queries.ContainerRepositoryItem](ctx, viper.GetViper())
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

	out, err := q.GetConfig(ctx, organisation.Id, configCode)
	if err != nil {
		cfg.WriteStderr("failed to get config")
		return err
	}

	cfg.WriteList(out.ContainerRepositories)
	return nil
}
