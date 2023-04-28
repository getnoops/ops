package cmd

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the configuration",
	Long:  `init initializes the configuration.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithCancel(cmd.Context())
		defer cancel()

		fs := afero.NewOsFs()
		if err := fs.MkdirAll("./.bin", 0755); err != nil {
			return err
		}

		v := "1.3.7"
		installer := &releases.ExactVersion{
			InstallDir: "./.bin",
			Product:    product.Terraform,
			Version:    version.Must(version.NewVersion(v)),
		}

		execPath, err := installer.Install(ctx)
		if err != nil {
			return err
		}

		fmt.Printf("Installed %s to %s", v, execPath)

		cancel()
		<-ctx.Done()
		return nil
	},
}
