package login

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/getnoops/ops/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zitadel/oidc/v2/pkg/oidc"
	"go.uber.org/zap"
)

type Config struct {
}

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login to NoOps",
		Long:  `Using SSO login to NoOps`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			config, err := config.New[Config](ctx, viper.GetViper())
			if err != nil {
				return err
			}

			Login(config)
			return nil
		},
	}

	addFlags(cmd)
	return cmd
}

func Login(config *config.NoOps[Config]) {
	ctx, cancel := context.WithCancel(context.Background())
	tokenChan := make(chan *oidc.Tokens[*oidc.IDTokenClaims], 1)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	go func() {
		<-sigs
		cancel()
	}()

	server, err := NewServer(ctx, config, tokenChan)
	if err != nil {
		config.Log.Fatal("failed to create server", zap.Error(err))
	}

	select {
	case <-ctx.Done():
		os.Exit(0)
	case token := <-tokenChan:
		if err := server.Shutdown(ctx); err != nil {
			config.Log.Fatal("failed to shutdown server", zap.Error(err))
		}

		if err := config.StoreToken(token); err != nil {
			config.Log.Fatal("failed to store token", zap.Error(err))
		}
	}
}
