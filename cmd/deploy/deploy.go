package deploy

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/avast/retry-go/v4"
	"github.com/getnoops/ops/pkg/brain"
	"github.com/getnoops/ops/pkg/poller"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func New(brain brain.Manager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy a stack file",
		Long:  `Deploy a stack file to the specified environment.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			config := MustNewConfig(viper.GetViper())

			return Deploy(context.Background(), config, brain)
		},
	}

	addFlags(cmd)
	return cmd
}

func Deploy(ctx context.Context, config *Config, b brain.Manager) error {
	// Make sure that stack file path from flag actually exists in user directory
	if _, err := os.Stat(config.StackFile); err != nil {
		return err
	}

	newDeployment, err := b.CreateDeployment(ctx, config.Environment)
	if err != nil {
		return err
	}

	err = UploadStackFileToS3WithRetry(ctx, b, newDeployment.DeploymentId, config.StackFile, newDeployment.UploadUrl)
	if err != nil {
		return err
	}

	fmt.Println("Stack file uploaded.")

	err = b.NotifyStackFileUploadCompleted(ctx, newDeployment.DeploymentId, brain.NotifyUploadCompleteRequest{Success: true})
	if err != nil {
		return err
	}

	fmt.Println("Brain notified of stack file upload.")

	poller.Wait(ctx, poller.WaitOptions{
		DeploymentId: newDeployment.DeploymentId,
		ExecToken:    &newDeployment.SessionToken,
		PollerConfig: poller.PollConfig{Interval: 10, Expiry: 60},
		BrainClient:  b,
	})

	return nil
}

func UploadStackFileToS3(stackFilePath, uploadUrl string) error {
	file, err := os.Open(stackFilePath)
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

func UploadStackFileToS3WithRetry(ctx context.Context, b brain.Manager, deploymentId, stackFilePath, uploadUrl string) error {
	err := retry.Do(
		func() error {
			err := UploadStackFileToS3(stackFilePath, uploadUrl)
			return err
		},
		retry.Attempts(3),
		retry.OnRetry(func(n uint, err error) {
			log.Printf("Unable to upload Stack file to S3 bucket. Retrying request after error: %v", err)
		}),
	)
	if err != nil {
		e := err.Error()
		notifyStackFileUploadBody := brain.NotifyUploadCompleteRequest{Success: false, Error: &e}
		b.NotifyStackFileUploadCompleted(ctx, deploymentId, notifyStackFileUploadBody)
		return err
	}

	return nil
}
