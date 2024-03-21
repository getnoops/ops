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

type KeyResult struct {
	Id    uuid.UUID `json:"id"`
	Token string    `json:"token"`
}

type CreateConfig struct {
}

func CreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "create",
		Short:  "Will create an api key for a given compute, storage or integration",
		PreRun: util.BindPreRun,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			return Create(ctx)
		},
	}
	return cmd
}

func Create(ctx context.Context) error {
	cfg, err := config.New[CreateConfig, KeyResult](ctx, viper.GetViper())
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

	out, err := q.CreateApiKey(ctx, organisation.Id)
	if err != nil {
		cfg.WriteStderr("failed to create api key")
		return err
	}

	result := KeyResult{
		Id:    out.Id,
		Token: out.Token,
	}

	cfg.WriteObject(result)
	return nil
}
