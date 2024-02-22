package containerrepository

import (
	"context"

	"github.com/getnoops/ops/pkg/config"
	"github.com/getnoops/ops/pkg/queries"
	"github.com/getnoops/ops/pkg/util"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type CreateConfig struct {
}

func CreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "create [compute] [code]",
		Short:  "Will create a container repository for a given compute",
		Args:   cobra.ExactArgs(2),
		PreRun: util.BindPreRun,
		RunE: func(cmd *cobra.Command, args []string) error {
			configCode := args[0]
			code := args[1]

			ctx := cmd.Context()
			return Create(ctx, configCode, code)
		},
		ValidArgs: []string{"compute", "code"},
	}
	return cmd
}

func Create(ctx context.Context, computeCode string, code string) error {
	cfg, err := config.New[CreateConfig, *uuid.UUID](ctx, viper.GetViper())
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

	config, err := q.GetConfig(ctx, organisation.Id, computeCode)
	if err != nil {
		cfg.WriteStderr("failed to get configs")
		return nil
	}

	id := uuid.New()
	out, err := q.CreateContainerRepository(ctx, organisation.Id, id, config.Id, code)
	if err != nil {
		cfg.WriteStderr("failed to create container repository")
		return nil
	}

	cfg.WriteObject(out)
	return nil
}
