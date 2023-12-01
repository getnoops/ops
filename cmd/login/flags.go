package login

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func addFlags(cmd *cobra.Command) {
	bindStringFlag(cmd, "company", "The company")
}

func bindStringFlag(cmd *cobra.Command, name, description string) {
	cmd.PersistentFlags().String(name, viper.GetString(name), description)
	viper.BindPFlag(name, cmd.PersistentFlags().Lookup(name))
}
