package containerregistry

import (
	"context"
	"encoding/json"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/getnoops/ops/pkg/config"
	"github.com/getnoops/ops/pkg/queries"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type CreateConfig struct {
}

func CreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [compute] [code]",
		Short: "Will create a container registry for a given compute",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			configCode := args[0]
			code := args[1]

			ctx := cmd.Context()
			return Create(ctx, configCode, code)
		},
		ValidArgs: []string{"compute", "code"},
	}
	return cmd
}

func Create(ctx context.Context, computeCode string, code string) error {
	cfg, err := config.New[CreateConfig](ctx, viper.GetViper())
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

	out, err := q.CreateContainerRegistry(ctx, organisation.Id, config.Id, code)
	if err != nil {
		cfg.WriteStderr("failed to create container registry")
		return nil
	}

	switch cfg.Global.Format {
	case "table":
		t := table.New().
			Border(lipgloss.NormalBorder()).
			BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
			Headers("Id", "Config Code", "Config Name", "Code")

		t.Row(out.String(), config.Code, config.Name, code)

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
