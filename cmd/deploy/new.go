package configs

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/getnoops/ops/pkg/config"
	"github.com/getnoops/ops/pkg/queries"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type DeployConfig struct {
}

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy [code] [env]",
		Short: "Get a config",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			env := args[0]
			code := args[1]
			versionNumber := args[2]

			ctx := cmd.Context()
			return Deploy(ctx, env, code, versionNumber)
		},
		ValidArgs: []string{"env", "code", "version_number"},
	}
	return cmd
}

func GetEnvironment(ctx context.Context, q queries.Queries, organisationId uuid.UUID, code string) (*queries.Environment, error) {
	paged, err := q.GetEnvironments(ctx, organisationId, []string{code}, 1, 999)
	if err != nil {
		return nil, err
	}

	for _, env := range paged.Items {
		if env.Code == code {
			return &env, nil
		}
	}

	return nil, fmt.Errorf("environment not found")
}

func Deploy(ctx context.Context, env string, code string, versionNumber string) error {
	cfg, err := config.New[DeployConfig](ctx, viper.GetViper())
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
	environment, err := GetEnvironment(ctx, q, organisation.Id, env)
	if err != nil {
		cfg.WriteStderr("environment not found for config")
		return nil
	}

	// out, err := q.NewDeployment(ctx, organisation.Id, config.Id, environment.Id, configRevisionId, revisionId)
	// if err != nil {
	// 	cfg.WriteStderr("failed to deploy")
	// 	return nil
	// }

	switch cfg.Global.Format {
	case "table":
		t := table.New().
			Border(lipgloss.NormalBorder()).
			BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
			Headers("Environment Code", "Environment Name", "Config Code", "Config Name", "Config Version")

		t.Row(environment.Code, environment.Name, config.Code, config.Name, config.Version_number)

		cfg.WriteStdout(t.Render())
	case "json":
		out, _ := json.Marshal(config)
		cfg.WriteStdout(string(out))
	case "yaml":
		out, _ := yaml.Marshal(config)
		cfg.WriteStdout(string(out))
	}
	return nil
}
