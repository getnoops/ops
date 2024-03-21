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

type DeleteConfig struct {
}

func DeleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "delete [compute] [code]",
		Short:  "Will delete a container repository for a given compute",
		Args:   cobra.ExactArgs(2),
		PreRun: util.BindPreRun,
		RunE: func(cmd *cobra.Command, args []string) error {
			configCode := args[0]
			code := args[1]

			ctx := cmd.Context()
			return Delete(ctx, configCode, code)
		},
		ValidArgs: []string{"compute", "code"},
	}
	return cmd
}

func Delete(ctx context.Context, computeCode string, code string) error {
	cfg, err := config.New[DeleteConfig, *uuid.UUID](ctx, viper.GetViper())
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
		return err
	}

	containerRepository, err := GetRepository(config.ContainerRepositories, code)
	if err != nil {
		return err
	}

	out, err := q.DeleteContainerRepository(ctx, organisation.Id, containerRepository.Id)
	if err != nil {
		return err
	}

	//		t.Row(out.String(), config.Code, config.Name, code)

	cfg.WriteObject(out)
	return nil
}
