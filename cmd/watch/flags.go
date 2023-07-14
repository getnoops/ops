package watch

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func addFlags(cmd *cobra.Command) {
	bindStringPFlag(cmd, "deployment", "d", "", "The Deployment ID you want to watch", true)
}

func bindStringPFlag(cmd *cobra.Command, name, shorthand, defaultValue, description string, required bool) {
	cmd.PersistentFlags().StringP(name, shorthand, defaultValue, description)
	viper.BindPFlag(name, cmd.PersistentFlags().Lookup(name))

	if required {
		cmd.MarkPersistentFlagRequired(name)
	}
}
