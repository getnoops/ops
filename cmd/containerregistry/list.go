package containerregistry

import (
	"context"

	"github.com/getnoops/ops/pkg/config"
	"github.com/getnoops/ops/pkg/queries"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

type ListConfig struct {
}

func ListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [compute]",
		Short: "list container registries",
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

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		Headers("Code", "State")

	for _, item := range out.ContainerRegistries {
		t.Row(item.Code, string(item.State))
	}

	cfg.WriteStdout(t.Render())
	return nil
}
