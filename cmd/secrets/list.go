package secrets

import (
	"context"

	"github.com/getnoops/ops/pkg/config"
	"github.com/getnoops/ops/pkg/models"
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
		Short:  "list secrets for a given compute",
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
	cfg, err := config.New[ListConfig, *models.Secret](ctx, viper.GetViper())
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

	var secrets []*models.Secret
	for _, secret := range out.Secrets {
		secrets = append(secrets, models.ToSecret(secret))
	}

	cfg.WriteList(secrets)
	return nil
}
