package configs

import (
	"github.com/spf13/cobra"
)

type Config struct {
}

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "configs",
		Short: "Config commands",
	}

	cmd.AddCommand(ListCommand())
	return cmd
}
