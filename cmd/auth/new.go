package auth

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login to NoOps",
		Long:  `Using SSO login to NoOps`,
		Run: func(cmd *cobra.Command, args []string) {
			config := MustNewConfig(viper.GetViper())

			Login(config)
		},
	}

	addFlags(cmd)
	return cmd
}

func Login(config *Config) {
}
