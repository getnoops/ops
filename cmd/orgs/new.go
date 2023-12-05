package orgs

import (
	"github.com/spf13/cobra"
)

type Config struct {
}

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "orgs",
		Short: "Orgs commands",
	}

	cmd.AddCommand(ListCommand())
	return cmd
}
