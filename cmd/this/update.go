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

type UpdateConfig struct {
	Next     bool     `mapstructure:"next" default:""`
	File     string   `mapstructure:"file" default:"noops.yaml"`
	VarFiles []string `mapstructure:"var-file" default:"noops.yaml"`
	Deploy   string   `mapstructure:"deploy" default:""`
}

func UpdateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "update",
		Short:  "Use the noops file to update the configuration",
		PreRun: util.BindPreRun,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			return Upgrade(ctx)
		},
	}

	util.BindStringPFlag(cmd, "file", "f", "The yaml file with the configuration", "")
	util.BindBoolFlag(cmd, "next", "Use the next minor version", false)
	util.BindStringFlag(cmd, "deploy", "Deploy the configuration to environment", "")
	util.BindStringSliceFlag(cmd, "var-file", "Environment like files to update the noops file", []string{})
	return cmd
}

func Upgrade(ctx context.Context) error {
	cfg, err := config.New[UpdateConfig, *models.Config](ctx, viper.GetViper())
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

	hasDeploy := cfg.Command.Deploy != ""
	var environmentId uuid.UUID

	if hasDeploy {
		codes := []string{cfg.Command.Deploy}
		paged, err := q.GetEnvironments(ctx, organisation.Id, codes, 1, 1)
		if err != nil {
			cfg.WriteStderr("failed to find environment")
			return err
		}
		if len(paged.Items) == 0 {
			cfg.WriteStderr("environment not found")
			return nil
		}
		environmentId = paged.Items[0].Id
	}

	rev, err := models.LoadFile[models.NoOpsConfig](cfg.Command.File, models.WithOsEnv(), models.WithVarFiles(cfg.Command.VarFiles))
	if err != nil {
		cfg.WriteStderr("failed to read file")
		return err
	}
	if err := rev.Validate(); err != nil {
		cfg.WriteStderr("failed to validate file")
		return err
	}

	resourceInput := []*queries.ResourceInput{}
	for _, resource := range rev.Resources {
		resourceInput = append(resourceInput, &queries.ResourceInput{
			Code: resource.Code,
			Type: resource.Type,
			Data: resource.Data,
		})
	}
	revId := uuid.New()

	config, err := q.GetConfig(ctx, organisation.Id, rev.Code)
	if err != nil {
		cfg.WriteStderr("failed to get configs")
		return nil
	}

	updateConfig, err := q.UpdateConfig(ctx, &queries.UpdateConfigInput{
		Organisation_id: organisation.Id,
		Aggregate_id:    config.Id,
		Name:            config.Name,
		Resources:       resourceInput,
		Version_number:  config.Version_number,
		Revision_id:     revId,
	})
	if err != nil {
		cfg.WriteStderr("failed to update config")
		return err
	}

	cfg.WriteStdout(fmt.Sprintf("Updated config %s", updateConfig.String()))

	if hasDeploy {
		deploymentRevId := uuid.New()
		deployResult, err := q.NewDeployment(ctx, organisation.Id, environmentId, config.Id, revId, deploymentRevId)
		if err != nil {
			cfg.WriteStderr("failed to deploy")
			return err
		}

		cfg.WriteStdout(fmt.Sprintf("Deployed %s", deployResult.String()))
	}

	out := models.ToConfig(config)
	cfg.WriteObject(out)
	return nil
}
