package containerrepository

import (
	"context"

	"github.com/getnoops/ops/pkg/config"
	"github.com/getnoops/ops/pkg/queries"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type LoginConfig struct {
}

func LoginCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login [compute] [code]",
		Short: "Will return the password from ecr login for a given compute and code",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			configCode := args[0]
			code := args[1]

			ctx := cmd.Context()
			return Login(ctx, configCode, code)
		},
		ValidArgs: []string{"compute", "code"},
	}
	return cmd
}

func Login(ctx context.Context, computeCode string, code string) error {
	cfg, err := config.New[LoginConfig](ctx, viper.GetViper())
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

	repository, err := GetRepository(config.ContainerRepositories, code)
	if err != nil {
		cfg.WriteStderr("failed to get container repository")
		return nil
	}

	out, err := q.LoginContainerRepository(ctx, organisation.Id, repository.Id)
	if err != nil {
		cfg.WriteStderr("failed to authenticate container registry")
		return nil
	}

	cfg.WriteStdout(out.Password)
	return nil
}
