package util

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func BindStringFlag(cmd *cobra.Command, name string, description string, value string) {
	cmd.Flags().String(name, value, description)
}

func BindStringPFlag(cmd *cobra.Command, name string, shorthand string, description string, value string) {
	cmd.Flags().StringP(name, shorthand, value, description)
}

func BindStringSliceFlag(cmd *cobra.Command, name string, description string, value []string) {
	cmd.Flags().StringSlice(name, value, description)
}

func BindBoolFlag(cmd *cobra.Command, name string, description string, value bool) {
	cmd.Flags().Bool(name, value, description)
}

func BindIntFlag(cmd *cobra.Command, name string, description string, value int) {
	cmd.Flags().Int(name, value, description)
}

func BindPreRun(cmd *cobra.Command, args []string) {
	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		viper.BindPFlag("command."+flag.Name, flag)
	})
}

func BindStringPersistentFlag(cmd *cobra.Command, name, description string, value string) {
	cmd.PersistentFlags().String(name, value, description)
	flag := cmd.PersistentFlags().Lookup(name)

	viper.BindPFlag("global."+name, flag)
}
