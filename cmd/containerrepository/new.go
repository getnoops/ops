package containerrepository

import (
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "container-registry",
		Short: "Container Registry commands",
	}

	cmd.AddCommand(ListCommand())
	cmd.AddCommand(GetCommand())
	cmd.AddCommand(LoginCommand())
	cmd.AddCommand(CreateCommand())
	cmd.AddCommand(DeleteCommand())
	return cmd
}