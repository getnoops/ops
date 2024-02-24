package this

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
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
	Watch    bool     `mapstructure:"watch" default:"false"`
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
	util.BindBoolFlag(cmd, "watch", "Watch deployment for success", false)
	return cmd
}

func GetEnvironment(ctx context.Context, q queries.Queries, organisation *queries.Organisation, code string) (*queries.Environment, error) {
	if len(code) == 0 {
		return nil, nil
	}

	codes := []string{code}
	paged, err := q.GetEnvironments(ctx, organisation.Id, codes, 1, 1)
	if err != nil {
		return nil, err
	}
	if len(paged.Items) == 0 {
		return nil, fmt.Errorf("environment not found")
	}
	return paged.Items[0], nil
}

func GetVersion(versionNumber string, next bool) (string, error) {
	v, err := semver.NewVersion(versionNumber)
	if err != nil {
		return "", err
	}

	if next {
		return v.IncMinor().String(), nil
	}
	return v.String(), nil
}

func Watch(ctx context.Context, cfg *config.NoOps[UpdateConfig, *models.Config], q queries.Queries, organisation *queries.Organisation, deploymentRevisionId uuid.UUID, count int) error {
	if count > 10 {
		cfg.WriteStderr("deployment state unknown")
		return fmt.Errorf("deployment state unknown")
	}

	revision, err := q.GetDeploymentRevision(ctx, organisation.Id, deploymentRevisionId)
	if err != nil {
		cfg.WriteStderr("failed to get deployment")
		return err
	}

	asString := string(revision.State)
	if strings.HasSuffix(asString, "ing") {
		cfg.WriteStdout(fmt.Sprintf("Deployment still %s, waiting 30s\n", asString))
		time.Sleep(30 * time.Second)
		return Watch(ctx, cfg, q, organisation, deploymentRevisionId, count+1)
	}

	cfg.WriteStdout(fmt.Sprintf("Deployment %s\n", asString))

	if revision.State == queries.StackStateFailed {
		return fmt.Errorf("deployment failed")
	}
	return nil
}

func Deploy(ctx context.Context, cfg *config.NoOps[UpdateConfig, *models.Config], q queries.Queries, organisation *queries.Organisation, environment *queries.Environment, config *queries.Config, configRevisionId uuid.UUID, watch bool) error {
	if environment == nil {
		return nil
	}

	deploymentRevisionId := uuid.New()
	_, err := q.NewDeployment(ctx, organisation.Id, environment.Id, config.Id, configRevisionId, deploymentRevisionId)
	if err != nil {
		return err
	}

	cfg.WriteStdout(fmt.Sprintf("Deploying %s to %s\n", config.Code, environment.Code))

	if watch {
		return Watch(ctx, cfg, q, organisation, deploymentRevisionId, 1)
	}
	return nil
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

	rev, err := models.LoadFile[models.NoOpsConfig](cfg.Command.File, models.WithOsEnv(), models.WithVarFiles(cfg.Command.VarFiles))
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

	resourceInput := []*queries.ResourceInput{}
	for _, resource := range rev.Resources {
		resourceInput = append(resourceInput, &queries.ResourceInput{
			Code: resource.Code,
			Type: resource.Type,
			Data: resource.Data,
		})
	}

	versionNumber, err := GetVersion(config.Version_number, cfg.Command.Next)
	if err != nil {
		cfg.WriteStderr("failed to get version")
		return err
	}

	environment, err := GetEnvironment(ctx, q, organisation, cfg.Command.Deploy)
	if err != nil {
		cfg.WriteStderr("failed to get environment")
		return err
	}

	revId := uuid.New()
	if _, err := q.UpdateConfig(ctx, &queries.UpdateConfigInput{
		Organisation_id: organisation.Id,
		Aggregate_id:    config.Id,
		Name:            config.Name,
		Resources:       resourceInput,
		Version_number:  versionNumber,
		Revision_id:     revId,
		Access:          rev.Access,
	}); err != nil {
		cfg.WriteStderr("failed to update config")
		return err
	}

	cfg.WriteStdout(fmt.Sprintf("Updated config %s %s\n", config.Code, versionNumber))

	return Deploy(ctx, cfg, q, organisation, environment, config, revId, cfg.Command.Watch)
}
