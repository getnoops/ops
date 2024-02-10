package keys

import (
	"context"
	"encoding/json"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/contextcloud/goutils/xstring"
	"github.com/getnoops/ops/pkg/config"
	"github.com/getnoops/ops/pkg/queries"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type CreateConfig struct {
}

func CreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [compute|storage|integration] [code]",
		Short: "Will create an api key for a given compute, storage or integration",
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

	key, err := xstring.GenerateString(22, 4, 4, false, false)
	if err != nil {
		cfg.WriteStderr("failed to generate key")
		return nil
	}

	out, err := q.CreateApiKey(ctx, organisation.Id, config.Id, code, key)
	if err != nil {
		cfg.WriteStderr("failed to create api key")
		return nil
	}

	result := struct {
		Id   uuid.UUID `json:"id"`
		Code string    `json:"code"`
		Key  string    `json:"key"`
	}{
		Id:   out,
		Code: code,
		Key:  key,
	}

	switch cfg.Global.Format {
	case "table":
		t := table.New().
			Border(lipgloss.NormalBorder()).
			BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
			Headers("Id", "Code", "Key")

		t.Row(result.Id.String(), result.Code, result.Key)

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
