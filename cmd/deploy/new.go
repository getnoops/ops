package deploy

import (
	"github.com/spf13/cobra"
)

type Config struct {
}

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy commands",
	}

	cmd.AddCommand(ApplyCommand())
	return cmd
}
