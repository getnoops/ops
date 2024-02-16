package settings

import (
	"context"

	"github.com/getnoops/ops/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ViewConfig struct {
}

func ViewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "view",
		Short: "view all settings",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			return View(ctx)
		},
	}

	return cmd
}

func View(ctx context.Context) error {
	cfg, err := config.New[ViewConfig, map[string]string](ctx, viper.GetViper())
	if err != nil {
		return err
	}

	settings, err := cfg.GetSettings()
	if err != nil {
		return err
	}

	cfg.WriteObject(settings)
	return nil
}
