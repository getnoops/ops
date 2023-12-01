package upgrade

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/getnoops/ops/pkg/config"
	"github.com/getnoops/ops/pkg/selfupdate"
	"github.com/getnoops/ops/pkg/version"
)

type Config struct {
	Prerelease bool `default:"false"`
	Draft      bool `default:"false"`
}

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrades ops tool to the latest version",
		Long:  `Upgrade will check for the latest version and upgrade if necessary.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			config, err := config.New[Config](ctx, viper.GetViper())
			if err != nil {
				return err
			}
			return Update(config)
		},
	}

	addFlags(cmd)
	return cmd
}

func Update(config *config.NoOps[Config]) error {
	ctx := context.Background()

	updater, err := selfupdate.NewUpdater("getnoops/ops", config.Command.Prerelease, config.Command.Draft)
	if err != nil {
		return fmt.Errorf("error occurred while creating updater: %w", err)
	}

	latest, err := updater.GetLatest(ctx)
	if err != nil {
		return fmt.Errorf("error occurred while detecting version: %w", err)
	}

	commit := version.Commit()
	diff, err := selfupdate.IsDifferent(commit, latest.Filename)
	if err != nil {
		return fmt.Errorf("error occurred while checking for latest version: %w", err)
	}

	if !diff {
		log.Println("You already have the latest")
		return nil
	}

	exePath, err := os.Executable()
	if err != nil {
		return errors.New("could not locate executable path")
	}
	if err := updater.UpdateTo(ctx, latest, exePath); err != nil {
		return fmt.Errorf("error occurred while updating binary: %w", err)
	}
	log.Println("Successfully updated ops")
	return nil
}
