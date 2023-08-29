package list

import (
	"context"
	"fmt"

	"github.com/getnoops/ops/pkg/brain"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func New(brain brain.Manager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List active deployments.",
		Long:  "List all deployments that have a status of either `PENDING` or `RUNNING`.",
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = MustNewConfig(viper.GetViper())

			return ListActiveDeployments(brain)
		},
	}

	return cmd
}

func ListActiveDeployments(b brain.Manager) error {
	activeDeployments, err := b.ListActiveDeployments(context.Background())
	if err != nil {
		return err
	}

	for _, d := range *activeDeployments {
		fmt.Printf("\n - %s (%s): %s", d.Status, d.EnvironmentName, d.DeploymentId)
	}

	return nil
}
