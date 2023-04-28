package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"cli/pkg/parser"
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply the deployment descriptor",
	Long:  `apply applies the deployment descriptor.`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithCancel(cmd.Context())
		defer cancel()

		filename := args[0]
		file, err := os.Open(filename)
		if err != nil {
			return err
		}

		p := parser.NewParser()

		spec, err := p.Parse(ctx, file)
		if err != nil {
			return err
		}

		// do something?

		fmt.Printf("%+v\n", spec)
		return nil
	},
}
