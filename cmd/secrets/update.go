package secrets

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
		Use:    "update [compute] [env] [code] [value]",
		Short:  "Will update a secret for a given compute",
		Args:   cobra.ExactArgs(4),
		PreRun: util.BindPreRun,
		RunE: func(cmd *cobra.Command, args []string) error {
			configCode := args[0]
			environmentCode := args[1]
			code := args[2]
			value := args[3]

			ctx := cmd.Context()
			return Update(ctx, configCode, environmentCode, code, value)
		},
		ValidArgs: []string{"compute", "env", "code", "value"},
	}
	return cmd
}

func Update(ctx context.Context, computeCode string, environmentCode string, code string, value string) error {
	cfg, err := config.New[UpdateConfig, *uuid.UUID](ctx, viper.GetViper())
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
		return err
	}

	environment, err := GetEnvironment(ctx, q, organisation, environmentCode)
	if err != nil {
		cfg.WriteStderr("failed to get environment")
		return err
	}

	secret, err := GetSecret(config.Secrets, environmentCode, code)
	if err != nil {
		cfg.WriteStderr("secret not found")
		return err
	}

	out, err := q.UpdateSecret(ctx, organisation.Id, secret.Id, config.Id, environment.Id, code, value)
	if err != nil {
		cfg.WriteStderr("failed to create container repository")
		return err
	}

	cfg.WriteObject(out)
	return nil
}
