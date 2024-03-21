package this

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/getnoops/ops/pkg/config"
	"github.com/getnoops/ops/pkg/models"
	"github.com/getnoops/ops/pkg/queries"
	"github.com/getnoops/ops/pkg/util"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type DeleteConfig struct {
	File  string `mapstructure:"file" default:"noops.yaml"`
	Watch bool   `mapstructure:"watch" default:"false"`
}

func DeleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "delete [env]",
		Short:  "Use the noops file to delete the configuration",
		Args:   cobra.ExactArgs(1),
		PreRun: util.BindPreRun,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			environmentCode := args[0]

			return Delete(ctx, environmentCode)
		},
	}

	util.BindStringPFlag(cmd, "file", "f", "The yaml file with the configuration", "")
	util.BindBoolFlag(cmd, "watch", "Watch deployment for success", false)
	return cmd
}

func WatchDeployment(ctx context.Context, cfg *config.NoOps[DeleteConfig, *models.Config], q queries.Queries, organisation *queries.Organisation, deploymentId uuid.UUID) error {
	deployment, err := q.GetDeployment(ctx, organisation.Id, deploymentId)
	if err != nil {
		cfg.WriteStderr("failed to get deployment")
		return err
	}

	asString := string(deployment.State)
	if strings.HasSuffix(asString, "ing") {
		cfg.WriteStdout(fmt.Sprintf("Deployment still %s, waiting 30s", asString))
		time.Sleep(30 * time.Second)
		return WatchDeployment(ctx, cfg, q, organisation, deploymentId)
	}

	cfg.WriteStdout(fmt.Sprintf("Deployment %s", asString))

	if deployment.State == queries.StackStateFailed {
		return fmt.Errorf("deployment failed")
	}
	return nil
}

func Delete(ctx context.Context, environmentCode string) error {
	cfg, err := config.New[DeleteConfig, *models.Config](ctx, viper.GetViper())
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

	config, err := q.GetConfig(ctx, organisation.Id, rev.Code)
	if err != nil {
		cfg.WriteStderr("failed to get configs")
		return nil
	}

	environment, err := GetEnvironment(ctx, q, organisation, environmentCode)
	if err != nil {
		cfg.WriteStderr("failed to get environment")
		return err
	}

	deployment := GetDeployment(ctx, config, environment)
	if deployment == nil {
		cfg.WriteStderr("no deployment for environment")
		return err
	}

	if _, err := q.DeleteDeployment(ctx, organisation.Id, deployment.Id); err != nil {
		cfg.WriteStderr("failed to delete config")
		return err
	}

	cfg.WriteStdout(fmt.Sprintf("Deleting %s from %s", config.Code, environment.Code))

	if cfg.Command.Watch {
		return WatchDeployment(ctx, cfg, q, organisation, deployment.Id)
	}
	return nil
}
