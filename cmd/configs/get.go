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

type GetConfig struct {
}

func GetCommand(class queries.ConfigClass) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [code]",
		Short: "Get a config",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			code := args[0]

			ctx := cmd.Context()
			return Get(ctx, class, code)
		},
		ValidArgs: []string{"code"},
	}
	return cmd
}

func Get(ctx context.Context, class queries.ConfigClass, code string) error {
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

	config, err := q.GetConfig(ctx, organisation.Id, code)
	if err != nil {
		cfg.WriteStderr("failed to get configs")
		return nil
	}

	if config == nil || config.Class != class {
		cfg.WriteStderr("config not found")
		return nil
	}

	switch cfg.Global.Format {
	case "table":
		t := table.New().
			Border(lipgloss.NormalBorder()).
			BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
			Headers("Code", "Name", "State", "Version", "Revisions", "Registries")

		revisions := util.JoinStrings(config.Revisions, func(r queries.ConfigWithRevisionsRevisionsConfigRevision) string {
			return r.Version_number
		}, ", ")
		registries := util.JoinStrings(config.ContainerRegistries, func(r queries.ConfigWithRevisionsContainerRegistriesContainerRegistry) string {
			return r.Code
		}, ", ")

		t.Row(config.Code, config.Name, string(config.State), config.Version_number, revisions, registries)

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
