package keys

import (
	"context"
	"encoding/json"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
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
		Use:   "create [compute|storage|integration]",
		Short: "Will create an api key for a given compute, storage or integration",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			configCode := args[0]

			ctx := cmd.Context()
			return Create(ctx, configCode)
		},
		ValidArgs: []string{"compute"},
	}
	return cmd
}

func Create(ctx context.Context, computeCode string) error {
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

	out, err := q.CreateApiKey(ctx, organisation.Id)
	if err != nil {
		cfg.WriteStderr("failed to create api key")
		return nil
	}

	result := struct {
		Id    uuid.UUID `json:"id"`
		Token string    `json:"token"`
	}{
		Id:    out.Id,
		Token: out.Token,
	}

	switch cfg.Global.Format {
	case "table":
		t := table.New().
			Border(lipgloss.NormalBorder()).
			BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
			Headers("Id", "Token")

		t.Row(result.Id.String(), result.Token)

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
