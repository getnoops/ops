package deploy

import (
	"context"
	"fmt"

	"github.com/getnoops/ops/pkg/config"
	"github.com/getnoops/ops/pkg/queries"
	"github.com/getnoops/ops/pkg/util"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ApplyConfig struct {
}

func ApplyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "apply [env] [compute|storage|integration] [version_number]",
		Short:  "Will deploy either a compute, storage or integration to an environment",
		Args:   cobra.ExactArgs(3),
		PreRun: util.BindPreRun,
		RunE: func(cmd *cobra.Command, args []string) error {
			env := args[0]
			code := args[1]
			versionNumber := args[2]

			ctx := cmd.Context()
			return Apply(ctx, env, code, versionNumber)
		},
		ValidArgs: []string{"env", "code", "version_number"},
	}
	return cmd
}

func GetEnvironment(ctx context.Context, q queries.Queries, organisation *queries.Organisation, code string) (*queries.Environment, error) {
	if len(code) == 0 {
		return nil, nil
	}

	codes := []string{code}
	states := []queries.StackState{queries.StackStateCreated}
	paged, err := q.GetEnvironments(ctx, organisation.Id, codes, states, 1, 1)
	if err != nil {
		return nil, err
	}
	if len(paged.Items) == 0 {
		return nil, fmt.Errorf("environment not found")
	}
	return paged.Items[0], nil
}

func GetConfigRevision(revisions []*queries.RevisionItem, versionNumber string) (*queries.RevisionItem, error) {
	for _, revision := range revisions {
		if revision.Version_number == versionNumber {
			return revision, nil
		}
	}

	return nil, fmt.Errorf("revision not found")
}

func GetDeploymentId(ctx context.Context, config *queries.Config, environment *queries.Environment) uuid.UUID {
	for _, deployment := range config.Deployments {
		if deployment.Environment.Id == environment.Id {
			return deployment.Id
		}
	}
	return uuid.New()
}

func Apply(ctx context.Context, env string, code string, versionNumber string) error {
	cfg, err := config.New[ApplyConfig, *uuid.UUID](ctx, viper.GetViper())
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

	config, err := q.GetConfig(ctx, organisation.Id, code)
	if err != nil {
		cfg.WriteStderr("failed to get configs")
		return nil
	}

	// get the correct environment.
	environment, err := GetEnvironment(ctx, q, organisation, env)
	if err != nil {
		cfg.WriteStderr("environment not found for config")
		return nil
	}

	deploymentId := GetDeploymentId(ctx, config, environment)

	// find the right revision.
	revision, err := GetConfigRevision(config.Revisions, versionNumber)
	if err != nil {
		cfg.WriteStderr(fmt.Sprintf("revision not found with version number %s", versionNumber))
		return nil
	}

	deploymentRevisionId := uuid.New()
	out, err := q.NewDeployment(ctx, organisation.Id, deploymentId, config.Id, environment.Id, revision.Id, deploymentRevisionId)
	if err != nil {
		cfg.WriteStderr("failed to deploy")
		return nil
	}

	cfg.WriteObject(out)
	return nil
}
