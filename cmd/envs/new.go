package envs

import (
	"github.com/spf13/cobra"
)

type Config struct {
}

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "envs",
		Short: "Environment commands",
	}

	cmd.AddCommand(ListCommand())
	return cmd
}
