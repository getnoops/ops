package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/getnoops/ops/pkg/util"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of ops",
	Long:  `All software has versions. This is NoOps'`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("version %s %s %s\n", util.Version(), runtime.GOOS, runtime.GOARCH)
	},
}
