package list

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/getnoops/ops/pkg/brain"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func New(brainClient *brain.ClientWithResponses) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List active deployments.",
		Long:  "List all deployments that have a status of either `PENDING` or `RUNNING`.",
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = MustNewConfig(viper.GetViper())

			return ListActiveDeployments(brainClient)
		},
	}

	return cmd
}

func ListActiveDeployments(brainClient *brain.ClientWithResponses) error {
	res, err := brainClient.ListActiveDeploymentsWithResponse(context.Background())
	if err != nil {
		return err
	}

	var activeDeploymentsResponse brain.ListDeploymentsResponse
	json.Unmarshal(res.Body, &activeDeploymentsResponse)

	for _, d := range activeDeploymentsResponse.Deployments {
		fmt.Printf("\n - %s (%s): %s", d.Status, d.EnvironmentName, d.DeploymentId)
	}

	return nil
}
