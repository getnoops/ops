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

type UpdateConfig struct {
}

func UpdateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update [compute|storage|integration] [code]",
		Short: "Will update an api key for a given compute, storage or integration",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			configCode := args[0]
			idStr := args[1]

			id, err := uuid.Parse(idStr)
			if err != nil {
				return err
			}

			ctx := cmd.Context()
			return Update(ctx, configCode, id)
		},
		ValidArgs: []string{"compute", "code"},
	}
	return cmd
}

func Update(ctx context.Context, computeCode string, id uuid.UUID) error {
	cfg, err := config.New[UpdateConfig](ctx, viper.GetViper())
	if err != nil {
		return err
	}

	q, err := queries.New(ctx, cfg)
	if err != nil {
		return err
	}

	_, orgErr := q.GetCurrentOrganisation(ctx)
	if orgErr == config.ErrNoOrganisation {
		cfg.WriteStderr("no organisation set")
		return nil
	}
	if orgErr != nil {
		return orgErr
	}

	out, err := q.UpdateApiKey(ctx, id)
	if err != nil {
		cfg.WriteStderr("failed to update api key")
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
