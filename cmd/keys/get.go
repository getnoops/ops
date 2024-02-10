package keys

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

type GetConfig struct {
}

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [compute|storage|integration] [code]",
		Short: "Will get an api key for a given compute, storage or integration",
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

	apiKey, err := GetApiKey(config.ApiKeys, code)
	if err != nil {
		cfg.WriteStderr("failed to get container repository")
		return nil
	}

	switch cfg.Global.Format {
	case "table":
		t := table.New().
			Border(lipgloss.NormalBorder()).
			BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
			Headers("Id", "State", "Code", "Created At", "Updated At", "Deleted At", "Authed At")

		t.Row(apiKey.Id.String(), string(apiKey.State), apiKey.Code, apiKey.Created_at.String(), apiKey.Updated_at.String(), apiKey.Deleted_at.String(), apiKey.Authed_at.String())

		cfg.WriteStdout(t.Render())
	case "json":
		out, _ := json.Marshal(apiKey)
		cfg.WriteStdout(string(out))
	case "yaml":
		out, _ := yaml.Marshal(apiKey)
		cfg.WriteStdout(string(out))
	}
	return nil
}

func GetApiKey(apiKeys []queries.ApiKeyItem, code string) (*queries.ApiKeyItem, error) {
	for _, item := range apiKeys {
		if item.Code == code {
			return &item, nil
		}
	}

	return nil, fmt.Errorf("api key not found")
}
