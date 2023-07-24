package deploy

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/getnoops/ops/pkg/brain"
	"github.com/getnoops/ops/pkg/poller"
	"github.com/getnoops/ops/pkg/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy a stack file",
		Long:  `Deploy a stack file to the specified environment.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			config := MustNewConfig(viper.GetViper())

			return Deploy(config)
		},
	}

	addFlags(cmd)
	return cmd
}

func Deploy(config *Config) error {
	// Make sure that stack file path from flag actually exists in user directory
	if _, err := os.Stat(config.StackFile); err != nil {
		return err
	}

	body := brain.CreateDeploymentRequest{EnvironmentName: config.Environment}
	r, err := util.MakeBodyReaderFromType(body)
	if err != nil {
		return err
	}

	res, err := brain.Client.CreateNewDeploymentWithBodyWithResponse(context.Background(), "application/json", r)
	if err != nil {
		return err
	}

	var newDeployment brain.CreateDeploymentResponse
	err = json.Unmarshal(res.Body, &newDeployment)
	if err != nil {
		return err
	}

	err = UploadStackFile(config.StackFile, newDeployment.UploadUrl)
	if err != nil {
		return err
	}

	fmt.Println("Stack file uploaded.")

	_, err = brain.Client.NotifyStackFileUploadCompleted(context.Background(), newDeployment.DeploymentId)
	if err != nil {
		return err
	}

	fmt.Println("Brain notified of stack file upload.")

	poller.Wait(context.Background(), poller.WaitOptions{
		DeploymentId: newDeployment.DeploymentId,
		ExecToken:    &newDeployment.SessionToken,
		PollerConfig: poller.PollConfig{Interval: 10, Expiry: 60},
	})

	return nil
}

func UploadStackFile(stack, uploadUrl string) error {
	file, err := os.Open(stack)
	if err != nil {
		return err
	}

	defer file.Close()

	// S3 requires the file in binary format
	fileContent, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	buffer := bytes.NewBuffer(fileContent)

	req, err := http.NewRequest(http.MethodPut, uploadUrl, buffer)
	if err != nil {
		return err
	}

	// Required fields from AWS
	req.Header.Set("Content-Type", "text/yaml")
	req.ContentLength = int64(len(fileContent))

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return errors.New("Unable to upload stack file, status code: " + res.Status)
	}

	return nil
}
