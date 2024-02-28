package configs

import (
	"strings"

	"github.com/getnoops/ops/pkg/queries"
	"github.com/spf13/cobra"
)

type Config struct {
}

func New(name string, classes ...queries.ConfigClass) *cobra.Command {
	code := strings.ToLower(name)
	short := name + " commands"

	cmd := &cobra.Command{
		Use:   code,
		Short: short,
	}

	cmd.AddCommand(ListCommand(classes))
	cmd.AddCommand(GetCommand(classes))
	return cmd
}
