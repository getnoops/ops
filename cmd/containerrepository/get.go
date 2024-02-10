package containerrepository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/contextcloud/goutils/xmap"
	"github.com/getnoops/ops/pkg/config"
	"github.com/getnoops/ops/pkg/queries"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type GetConfig struct {
}

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [compute] [code]",
		Short: "Will get a container registry for a given compute repository",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			configCode := args[0]
			code := args[1]

			ctx := cmd.Context()
			return Get(ctx, configCode, code)
		},
		ValidArgs: []string{"compute", "code"},
	}
	return cmd
}

func Get(ctx context.Context, computeCode string, code string) error {
	cfg, err := config.New[GetConfig](ctx, viper.GetViper())
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

	// we want the stack outputs.
	out, err := q.GetContainerRepository(ctx, organisation.Id, repository.Id)
	if err != nil {
		cfg.WriteStderr("failed to get container repository")
		return nil
	}
	outputs := map[string]string{}
	for _, output := range out.Stack.Outputs {
		outputs[output.Output_key] = output.Output_value
	}

	result := struct {
		RepositoryUri  string `json:"repository_uri"`
		RepositoryName string `json:"repository_name"`
		Username       string `json:"username"`
	}{
		Username: "AWS",
	}
	if err := xmap.Decode(outputs, &result); err != nil {
		cfg.WriteStderr("failed to decode stack outputs")
		return nil
	}

	switch cfg.Global.Format {
	case "table":
		t := table.New().
			Border(lipgloss.NormalBorder()).
			BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
			Headers("Registry", "Repository", "Username")

		t.Row(result.RepositoryUri, result.RepositoryName, result.Username)

		cfg.WriteStdout(t.Render())
	case "json":
		out, _ := json.Marshal(result)
		cfg.WriteStdout(string(out))
	case "yaml":
		out, _ := yaml.Marshal(result)
		cfg.WriteStdout(string(out))
	}
	return nil
}

func GetRepository(repositories []queries.ContainerRepositoryItem, code string) (*queries.ContainerRepositoryItem, error) {
	for _, repository := range repositories {
		if repository.Code == code {
			return &repository, nil
		}
	}

	return nil, fmt.Errorf("repository not found")
}
