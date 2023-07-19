package watch

import (
	"context"
	"fmt"

	"github.com/getnoops/ops/pkg/brain"
	"github.com/getnoops/ops/pkg/poller"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func New(brainClient *brain.ClientWithResponses) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "watch",
		Short: "Watch a deployment via polling.",
		Long:  "Watch a specific deployment's events via polling.",
		RunE: func(cmd *cobra.Command, args []string) error {
			config := MustNewConfig(viper.GetViper())

			fmt.Printf("\nNow watching deployment: %s \n", config.DeploymentId)

			return poller.Wait(context.Background(), poller.WaitOptions{
				DeploymentId: config.DeploymentId,
				BrainClient:  brainClient,
			})
		},
	}

	addFlags(cmd)
	return cmd
}