package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply the stack file",
	Long:  `apply applies the stack file.`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("hello %v\n", args)
		return nil
	},
}
