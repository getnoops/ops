package keys

import (
	"context"
	"encoding/json"

	"github.com/getnoops/ops/pkg/config"
	"github.com/getnoops/ops/pkg/queries"
	"github.com/getnoops/ops/pkg/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

type ListConfig struct {
	Page     int `mapstructure:"page" default:"1"`
	PageSize int `mapstructure:"page-size" default:"10"`
}

func ListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list api keys",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			return List(ctx)
		},
	}

	util.BindIntFlag(cmd, "page", "The page to load", 1)
	util.BindIntFlag(cmd, "page-size", "The number of items in the page", 10)
	return cmd
}

func List(ctx context.Context) error {
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

	out, err := q.GetApiKeys(ctx, organisation.Id, cfg.Command.Page, cfg.Command.PageSize)
	if err != nil {
		cfg.WriteStderr("failed to get config")
		return err
	}

	switch cfg.Global.Format {
	case "table":
		t := table.New().
			Border(lipgloss.NormalBorder()).
			BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
			Headers("Id", "State", "Created At", "Updated At", "Deleted At", "Authed At")

		for _, item := range out.Items {
			t.Row(item.Id.String(), string(item.State), item.Created_at.String(), item.Updated_at.String(), item.Deleted_at.String(), item.Authed_at.String())
		}

		cfg.WriteStdout(t.Render())
	case "json":
		out, _ := json.Marshal(out.Items)
		cfg.WriteStdout(string(out))
	case "yaml":
		out, _ := yaml.Marshal(out.Items)
		cfg.WriteStdout(string(out))
	}
	return nil
}
