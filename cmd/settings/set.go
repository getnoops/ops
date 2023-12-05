package settings

import (
	"context"
	"fmt"
	"strings"

	"github.com/getnoops/ops/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type SetConfig struct {
}

func SetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set",
		Short: "set a No_Ops cli property",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			config, err := config.New[SetConfig](ctx, viper.GetViper())
			if err != nil {
				return err
			}

			return Set(ctx, config, args[0], args[1])
		},
	}

	return cmd
}

func Set(ctx context.Context, config *config.NoOps[SetConfig], key string, val string) error {
	settings, err := config.GetSettings()
	if err != nil {
		return err
	}

	switch strings.ToLower(key) {
	case "organisation":
		settings["organisation"] = val
	case "org":
		settings["organisation"] = val
	default:
		return fmt.Errorf("unknown setting %s, should be one of: [%s]", key, strings.Join(ValidProps, ","))
	}

	return config.StoreSettings(settings)
}
