package main

import (
	"os"

	"github.com/getnoops/ops/cmd"
	"github.com/spf13/cobra"
)

func main() {
	args := os.Args[1:]
	rootCmd := cmd.New(os.Stdout, os.Stdin, args)
	cobra.CheckErr(rootCmd.Execute())
}
