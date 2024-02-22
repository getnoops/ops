package this

import (
	"github.com/spf13/cobra"
)

type Config struct {
}

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "this",
		Short: "This commands",
	}

	cmd.AddCommand(InfoCommand())
	cmd.AddCommand(UpdateCommand())
	return cmd
}
