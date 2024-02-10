package keys

import (
	"context"
	"encoding/json"

	"github.com/getnoops/ops/pkg/config"
	"github.com/getnoops/ops/pkg/queries"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

type ListConfig struct {
}

func ListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [compute|storage|integration]",
		Short: "list api keys for a given compute, storage or integration",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			configCode := args[0]

			ctx := cmd.Context()
			return List(ctx, configCode)
		},
	}
	return cmd
}

func List(ctx context.Context, configCode string) error {
	cfg, err := config.New[ListConfig](ctx, viper.GetViper())
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

	out, err := q.GetConfig(ctx, organisation.Id, configCode)
	if err != nil {
		cfg.WriteStderr("failed to get config")
		return err
	}

	switch cfg.Global.Format {
	case "table":
		t := table.New().
			Border(lipgloss.NormalBorder()).
			BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
			Headers("Id", "State", "Code", "Created At", "Updated At", "Deleted At", "Authed At")

		for _, item := range out.ApiKeys {
			t.Row(item.Id.String(), string(item.State), item.Code, item.Created_at.String(), item.Updated_at.String(), item.Deleted_at.String(), item.Authed_at.String())
		}

		cfg.WriteStdout(t.Render())
	case "json":
		out, _ := json.Marshal(out.ApiKeys)
		cfg.WriteStdout(string(out))
	case "yaml":
		out, _ := yaml.Marshal(out.ApiKeys)
		cfg.WriteStdout(string(out))
	}
	return nil
}
