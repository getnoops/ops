package keys

import (
	"context"

	"github.com/getnoops/ops/pkg/config"
	"github.com/getnoops/ops/pkg/queries"
	"github.com/getnoops/ops/pkg/util"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type UpdateConfig struct {
}

func UpdateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "update [id]",
		Short:  "Will update an api key",
		Args:   cobra.ExactArgs(1),
		PreRun: util.BindPreRun,
		RunE: func(cmd *cobra.Command, args []string) error {
			idStr := args[0]

			id, err := uuid.Parse(idStr)
			if err != nil {
				return err
			}

			ctx := cmd.Context()
			return Update(ctx, id)
		},
		ValidArgs: []string{"id"},
	}
	return cmd
}

func Update(ctx context.Context, id uuid.UUID) error {
	cfg, err := config.New[UpdateConfig, KeyResult](ctx, viper.GetViper())
	if err != nil {
		return err
	}

	q, err := queries.New(ctx, cfg)
	if err != nil {
		return err
	}

	_, orgErr := q.GetCurrentOrganisation(ctx)
	if orgErr == config.ErrNoOrganisation {
		cfg.WriteStderr("no organisation set")
		return nil
	}
	if orgErr != nil {
		return orgErr
	}

	out, err := q.UpdateApiKey(ctx, id)
	if err != nil {
		cfg.WriteStderr("failed to update api key")
		return nil
	}

	result := KeyResult{
		Id:    out.Id,
		Token: out.Token,
	}

	cfg.WriteObject(result)
	return nil
}
