package upgrade

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func addFlags(cmd *cobra.Command) {
	bindBoolFlag(cmd, "prerelease", "Check for pre-release versions")
	bindBoolFlag(cmd, "draft", "Check for draft versions")

}

func bindBoolFlag(cmd *cobra.Command, name, description string) {
	cmd.PersistentFlags().Bool(name, viper.GetBool(name), description)
	viper.BindPFlag(name, cmd.PersistentFlags().Lookup(name))
}
