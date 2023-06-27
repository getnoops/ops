package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/getnoops/ops/pkg/selfupdate"
	"github.com/getnoops/ops/pkg/util"
)

func update(ctx context.Context, version string, prerelease bool, draft bool) error {
	updater, err := selfupdate.NewUpdater("getnoops/ops", prerelease, draft)
	if err != nil {
		return fmt.Errorf("error occurred while creating updater: %w", err)
	}

	latest, err := updater.GetLatest(ctx)
	if err != nil {
		return fmt.Errorf("error occurred while detecting version: %w", err)
	}

	fmt.Println(latest)

	// if latest.LessOrEqual(version) {
	// 	log.Printf("Current version (%s) is the latest", version)
	// 	return nil
	// }

	// exe, err := os.Executable()
	// if err != nil {
	// 	return errors.New("could not locate executable path")
	// }
	// if err := selfupdate.UpdateTo(context.Background(), latest.AssetURL, latest.AssetName, exe); err != nil {
	// 	return fmt.Errorf("error occurred while updating binary: %w", err)
	// }
	// log.Printf("Successfully updated to version %s", latest.Version())
	return nil
}

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrades ops tool to the latest version",
	Long:  `Upgrade will check for the latest version and upgrade if necessary.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return update(context.Background(), util.Version(), true, false)
	},
}
