package containerregistry

import (
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "container-registry",
		Short: "Container Registry commands",
	}

	cmd.AddCommand(ListCommand())
	cmd.AddCommand(CreateCommand())
	cmd.AddCommand(DeleteCommand())
	return cmd
}
