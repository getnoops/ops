package containerrepository

import (
	"context"

	"github.com/getnoops/ops/pkg/config"
	"github.com/getnoops/ops/pkg/queries"
	"github.com/getnoops/ops/pkg/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type LoginConfig struct {
}

func LoginCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "login",
		Short:  "Will return the password from ecr login",
		PreRun: util.BindPreRun,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			return Login(ctx)
		},
		ValidArgs: []string{"compute", "code"},
	}
	return cmd
}

func Login(ctx context.Context) error {
	cfg, err := config.New[LoginConfig, string](ctx, viper.GetViper())
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

	out, err := q.LoginContainerRepository(ctx, organisation.Id)
	if err != nil {
		cfg.WriteStderr("failed to authenticate container repository")
		return nil
	}

	cfg.WriteStdout(out.Password)
	return nil
}
