package secrets

import (
	"context"
	"fmt"

	"github.com/getnoops/ops/pkg/config"
	"github.com/getnoops/ops/pkg/queries"
	"github.com/getnoops/ops/pkg/util"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func GetEnvironment(ctx context.Context, q queries.Queries, organisation *queries.Organisation, code string) (*queries.Environment, error) {
	if len(code) == 0 {
		return nil, nil
	}

	codes := []string{code}
	states := []queries.StackState{queries.StackStateCreated}
	paged, err := q.GetEnvironments(ctx, organisation.Id, codes, states, 1, 1)
	if err != nil {
		return nil, err
	}
	if len(paged.Items) == 0 {
		return nil, fmt.Errorf("environment not found")
	}
	return paged.Items[0], nil
}

type CreateSecret struct {
}

func CreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "create [compute] [env] [code] [value]",
		Short:  "Will create a new secret for a given compute for a given environment",
		Args:   cobra.ExactArgs(4),
		PreRun: util.BindPreRun,
		RunE: func(cmd *cobra.Command, args []string) error {
			configCode := args[0]
			environmentCode := args[1]
			code := args[2]
			value := args[3]

			ctx := cmd.Context()
			return Create(ctx, configCode, environmentCode, code, value)
		},
		ValidArgs: []string{"compute", "env", "code", "value"},
	}
	return cmd
}

func Create(ctx context.Context, computeCode string, environmentCode string, code string, value string) error {
	cfg, err := config.New[CreateSecret, *uuid.UUID](ctx, viper.GetViper())
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

	environment, err := GetEnvironment(ctx, q, organisation, environmentCode)
	if err != nil {
		cfg.WriteStderr("failed to get environment")
		return nil
	}

	id := uuid.New()
	out, err := q.CreateSecret(ctx, organisation.Id, id, config.Id, environment.Id, code, value)
	if err != nil {
		cfg.WriteStderr("failed to create secret")
		return nil
	}

	cfg.WriteObject(out)
	return nil
}
