package this

import (
	"context"

	"github.com/getnoops/ops/pkg/config"
	"github.com/getnoops/ops/pkg/models"
	"github.com/getnoops/ops/pkg/queries"
	"github.com/getnoops/ops/pkg/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type InfoConfig struct {
	File string `mapstructure:"file" default:"noops.yaml"`
}

func InfoCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "info",
		Short:  "Use the noops file to get the information",
		PreRun: util.BindPreRun,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			return Info(ctx)
		},
	}

	util.BindStringPFlag(cmd, "file", "f", "The yaml file with the configuration", "")
	return cmd
}

func Info(ctx context.Context) error {
	cfg, err := config.New[InfoConfig, *models.Config](ctx, viper.GetViper())
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

	rev, err := models.LoadFile[models.NoOpsCode](cfg.Command.File)
	if err != nil {
		cfg.WriteStderr("failed to read file")
		return err
	}

	config, err := q.GetConfig(ctx, organisation.Id, rev.Code)
	if err != nil {
		cfg.WriteStderr("failed to get configs")
		return nil
	}

	out := models.ToConfig(config)
	cfg.WriteObject(out)
	return nil
}
