package settings

import (
	"context"
	"fmt"
	"strings"

	"github.com/getnoops/ops/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ValidProps = []string{"organisation", "org"}

type UnsetConfig struct {
}

func UnsetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unset",
		Short: "unset a No_Ops cli property",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			config, err := config.New[UnsetConfig](ctx, viper.GetViper())
			if err != nil {
				return err
			}

			return Unset(ctx, config, args[0])
		},
	}

	return cmd
}

func Unset[T any](ctx context.Context, config *config.NoOps[T], key string) error {
	settings, err := config.GetSettings()
	if err != nil {
		return err
	}

	switch strings.ToLower(key) {
	case "organisation":
		delete(settings, "organisation")
	case "org":
		delete(settings, "organisation")
	default:
		return fmt.Errorf("unknown setting %s, should be one of: [%s]", key, strings.Join(ValidProps, ","))
	}

	return config.StoreSettings(settings)
}
