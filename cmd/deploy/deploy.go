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

	newDeployment, err := CreateBrainDeployment(config.Environment)
	if err != nil {
		return err
	}

	fmt.Println("Deployment to brain created.")

	err = UploadStackFile(config.StackFile, newDeployment.UploadUrl)
	if err != nil {
		return err
	}

	fmt.Println("Stack file uploaded!")

	err = NotifyUploadComplete(newDeployment.DeploymentId)
	if err != nil {
		return err
	}

	fmt.Println("Brain notified of stack file upload.")

	poller.Wait(context.Background(), poller.WaitOptions{
		DeploymentId: newDeployment.DeploymentId,
		ExecToken:    newDeployment.SessionToken,
	})

	return nil
}

func CreateBrainDeployment(env string) (*brain.CreateDeploymentResponse, error) {
	url := viper.GetString("BrainUrl")
	body := brain.CreateDeploymentRequest{EnvironmentName: env}
	req, err := brain.NewCreateNewDeploymentRequest(url, body)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	resData, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var newDeployment brain.CreateDeploymentResponse
	json.Unmarshal(resData, &newDeployment)

	return &newDeployment, nil
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

func NotifyUploadComplete(deploymentId string) error {
	req, err := brain.NewNotifyUploadCompletedRequest(viper.GetString("BrainUrl"), deploymentId)
	if err != nil {
		return err
	}

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	return nil
}
