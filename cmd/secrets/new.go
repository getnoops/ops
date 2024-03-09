package secrets

import (
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "secrets",
		Short: "Secret commands",
	}

	cmd.AddCommand(ListCommand())
	cmd.AddCommand(GetCommand())
	cmd.AddCommand(CreateCommand())
	cmd.AddCommand(UpdateCommand())
	cmd.AddCommand(DeleteCommand())
	return cmd
}
