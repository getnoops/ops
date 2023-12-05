package util

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func BindStringFlag(cmd *cobra.Command, name, description string, value string) {
	cmd.Flags().String(name, value, description)
	viper.BindPFlag("Command."+name, cmd.Flags().Lookup(name))
}

func BindBoolFlag(cmd *cobra.Command, name, description string, value bool) {
	cmd.Flags().Bool(name, value, description)
	viper.BindPFlag("Command."+name, cmd.Flags().Lookup(name))
}

func BindIntFlag(cmd *cobra.Command, name, description string, value int) {
	cmd.Flags().Int(name, value, description)
	viper.BindPFlag("Command."+name, cmd.Flags().Lookup(name))
}

func BindStringPersistentFlag(cmd *cobra.Command, name, description string, value string) {
	cmd.PersistentFlags().String(name, value, description)
	viper.BindPFlag("Global."+name, cmd.PersistentFlags().Lookup(name))
}
