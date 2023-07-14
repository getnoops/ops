package list

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/getnoops/ops/pkg/brain"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List active deployments.",
		Long:  "List all deployments that have a status of either `PENDING` or `RUNNING`.",
		RunE: func(cmd *cobra.Command, args []string) error {
			config := MustNewConfig(viper.GetViper())

			return ListActiveDeployments(config)
		},
	}

	return cmd
}

func ListActiveDeployments(_ *Config) error {
	url := viper.GetString("BrainUrl")
	req, err := brain.NewListActiveDeploymentsRequest(url)
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	resData, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var activeDeploymentsResponse brain.ListDeploymentsResponse
	json.Unmarshal(resData, &activeDeploymentsResponse)

	for _, d := range activeDeploymentsResponse.Deployments {
		fmt.Printf("\n - %s (%s): %s", d.Status, d.EnvironmentName, d.DeploymentId)
	}

	return nil
}
