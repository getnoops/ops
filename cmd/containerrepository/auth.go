package containerrepository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/getnoops/ops/pkg/config"
	"github.com/getnoops/ops/pkg/queries"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type AuthConfig struct {
}

func AuthCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth [compute] [code]",
		Short: "Will authenticate a container registry for a given compute repository",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			configCode := args[0]
			code := args[1]

			ctx := cmd.Context()
			return Auth(ctx, configCode, code)
		},
		ValidArgs: []string{"compute", "code"},
	}
	return cmd
}

func Auth(ctx context.Context, computeCode string, code string) error {
	cfg, err := config.New[AuthConfig](ctx, viper.GetViper())
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

	switch cfg.Global.Format {
	case "table":
		t := table.New().
			Border(lipgloss.NormalBorder()).
			BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
			Headers("Registry", "Repository", "Username", "Password")

		t.Row(out.Registry_url, out.Repository_name, out.Username, out.Password)

		cfg.WriteStdout(t.Render())
	case "json":
		out, _ := json.Marshal(out)
		cfg.WriteStdout(string(out))
	case "yaml":
		out, _ := yaml.Marshal(out)
		cfg.WriteStdout(string(out))
	}
	return nil
}

func GetRepository(repositories []queries.ConfigWithRevisionsContainerRepositoriesContainerRepository, code string) (*queries.ConfigWithRevisionsContainerRepositoriesContainerRepository, error) {
	for _, repository := range repositories {
		if repository.Code == code {
			return &repository, nil
		}
	}

	return nil, fmt.Errorf("repository not found")
}
