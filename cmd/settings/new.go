package settings

import (
	"github.com/spf13/cobra"
)

type Config struct {
}

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "settings",
		Short: "Settings commands",
	}

	cmd.AddCommand(ViewCommand())
	cmd.AddCommand(SetCommand())
	cmd.AddCommand(UnsetCommand())
	return cmd
}
