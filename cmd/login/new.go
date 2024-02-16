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
			return Login(ctx)
		},
	}

	addFlags(cmd)
	return cmd
}

func Login(ctx context.Context) error {
	cfg, err := config.New[Config, string](ctx, viper.GetViper())
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(ctx)
	tokenChan := make(chan *oidc.Tokens[*oidc.IDTokenClaims], 1)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	go func() {
		<-sigs
		cancel()
	}()

	server, err := NewServer(ctx, cfg, tokenChan)
	if err != nil {
		cfg.WriteStderr("failed to create server")
		return err
	}

	select {
	case <-ctx.Done():
		os.Exit(0)
	case token := <-tokenChan:
		if err := server.Shutdown(ctx); err != nil {
			cfg.WriteStderr("failed to shutdown server")
			return err
		}

		if err := cfg.StoreToken(token); err != nil {
			cfg.WriteStderr("failed to store token")
			return err
		}
	}
	return nil
}
