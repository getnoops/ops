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

type RestoreConfig struct {
}

func RestoreCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "restore [compute] [env] [code]",
		Short:  "Will restore a secret for a given compute and environment. Must be in a pending deletion state.",
		Args:   cobra.ExactArgs(3),
		PreRun: util.BindPreRun,
		RunE: func(cmd *cobra.Command, args []string) error {
			configCode := args[0]
			envCode := args[1]
			code := args[2]

			ctx := cmd.Context()
			return Restore(ctx, configCode, envCode, code)
		},
		ValidArgs: []string{"compute", "code"},
	}
	return cmd
}

func Restore(ctx context.Context, computeCode string, environmentCode string, code string) error {
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
		cfg.WriteStderr("failed to get configs")
		return err
	}

	secret, err := GetSecret(config.Secrets, environmentCode, code)
	if err != nil {
		cfg.WriteStderr("secret not found")
		return err
	}

	out, err := q.DeleteSecret(ctx, organisation.Id, secret.Id)
	if err != nil {
		cfg.WriteStderr("failed to delete container repository")
		return err
	}

	cfg.WriteObject(out)
	return nil
}
