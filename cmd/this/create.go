package this

import (
	"context"
	"fmt"

	"github.com/getnoops/ops/pkg/config"
	"github.com/getnoops/ops/pkg/models"
	"github.com/getnoops/ops/pkg/queries"
	"github.com/getnoops/ops/pkg/util"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type CreateConfig struct {
	File string `mapstructure:"file" default:"noops.yaml"`
}

func CreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "create",
		Short:  "Use the noops file to create a configuration",
		PreRun: util.BindPreRun,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			return Create(ctx)
		},
	}

	util.BindStringPFlag(cmd, "file", "f", "The yaml file with the configuration", "")
	return cmd
}

func Create(ctx context.Context) error {
	cfg, err := config.New[CreateConfig, *models.Config](ctx, viper.GetViper())
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

	rev, err := models.LoadFile[models.NoOpsConfig](cfg.Command.File, models.WithOsEnv())
	if err != nil {
		cfg.WriteStderr("failed to read file")
		return err
	}
	if err := rev.Validate(); err != nil {
		cfg.WriteStderr("failed to validate file")
		return err
	}

	id := uuid.New()
	if _, err := q.CreateConfig(ctx, organisation.Id, id, rev.Name, rev.Code, rev.Class); err != nil {
		cfg.WriteStderr("failed to create config")
		return err
	}

	q.UpdateConfig(ctx, &queries.UpdateConfigInput{
		Aggregate_id:    id,
		Organisation_id: organisation.Id,
		Name:            rev.Name,
		Resources:       rev.Resources,
		Access:          rev.Access,
		Version_number:  "1.0.0",
		Revision_id:     uuid.New(),
	})

	cfg.WriteStdout(fmt.Sprintf("Created config %s(%s)", rev.Name, rev.Code))

	return nil
}
