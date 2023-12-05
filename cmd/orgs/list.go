package orgs

import (
	"context"

	"github.com/getnoops/ops/pkg/config"
	"github.com/getnoops/ops/pkg/queries"
	"github.com/getnoops/ops/pkg/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

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
		Short: "list projects accessible by the active account",
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

	queries, err := queries.New(ctx, cfg)
	if err != nil {
		return err
	}

	out, err := queries.GetMemberOrganisations(ctx, cfg.Command.Page, cfg.Command.PageSize)
	if err != nil {
		cfg.WriteStderr("failed to get member organisations")
		return err
	}

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		Headers("Code", "Name", "State")

	for _, item := range out.Items {
		t.Row(item.Code, item.Name, string(item.State))
	}

	cfg.WriteStdout(t.Render())
	return nil
}
