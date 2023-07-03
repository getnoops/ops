package deploy

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func addFlags(cmd *cobra.Command) {
	// Do we **only** want to allow `dev`, `test`, `stage` & `prod` values here, or can they be named anything?
	bindStringPFlag(cmd, "env", "e", "test", "The environment to deploy to", true)

	bindStringPFlag(cmd, "file", "f", "stack.yml", "The location of your stack file", true)
}

func bindStringPFlag(cmd *cobra.Command, name, shorthand, defaultValue, description string, required bool) {
	cmd.PersistentFlags().StringP(name, shorthand, defaultValue, description)
	viper.BindPFlag(name, cmd.PersistentFlags().Lookup(name))

	if required {
		cmd.MarkPersistentFlagRequired(name)
	}
}
