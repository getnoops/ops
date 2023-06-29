package cmd

import (
	"fmt"
	"net/http"

	"github.com/cli/oauth/device"
	"github.com/getnoops/ops/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to NoOps",
	Long:  `Using SSO login to NoOps`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.New()
		if err := viper.Unmarshal(cfg); err != nil {
			return err
		}

		requestUri := fmt.Sprintf("https://%s/login/device/code", cfg.Auth.Host)

		httpClient := http.DefaultClient
		code, err := device.RequestCode(httpClient, requestUri, cfg.Auth.ClientId, cfg.Auth.Scopes)
		if err != nil {
			return err
		}

		fmt.Printf("Copy code: %s\n", code.UserCode)
		fmt.Printf("then open: %s\n", code.VerificationURI)

		accessTokenUri := fmt.Sprintf("https://%s/login/oauth/access_token", cfg.Auth.Host)

		accessToken, err := device.PollToken(httpClient, accessTokenUri, cfg.Auth.ClientId, code)
		if err != nil {
			return err
		}

		fmt.Printf("Access token: %s\n", accessToken.Token)
		fmt.Printf("Refresh token: %s\n", accessToken.RefreshToken)
		return nil
	},
}
