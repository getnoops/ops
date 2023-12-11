package configs

import (
	"context"
	"encoding/json"

	"github.com/getnoops/ops/pkg/config"
	"github.com/getnoops/ops/pkg/queries"
	"github.com/getnoops/ops/pkg/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

type DescribeConfig struct {
}

func DescribeCommand(class queries.ConfigClass) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list projects accessible by the active account",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			return Describe(ctx, class)
		},
	}

	util.BindIntFlag(cmd, "page", "The page to load", 1)
	util.BindIntFlag(cmd, "page-size", "The number of items in the page", 10)
	return cmd
}

func Describe(ctx context.Context, class queries.ConfigClass) error {
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

	configs, err := q.GetConfigs(ctx, organisation.Id, class, cfg.Command.Page, cfg.Command.PageSize)
	if err != nil {
		cfg.WriteStderr("failed to get configs")
		return nil
	}

	switch cfg.Global.Format {
	case "table":
		t := table.New().
			Border(lipgloss.NormalBorder()).
			BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
			Headers("Code", "Name", "Status", "Version")

		for _, item := range configs.Items {
			t.Row(item.Code, item.Name, string(item.Status), item.Version_number)
		}

		cfg.WriteStdout(t.Render())
	case "json":
		out, _ := json.Marshal(configs)
		cfg.WriteStdout(string(out))
	case "yaml":
		out, _ := yaml.Marshal(configs)
		cfg.WriteStdout(string(out))
	}
	return nil
}
