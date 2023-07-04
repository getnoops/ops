package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// TODO: Auto generate these values from Brain API
type NewDeploymentRequest struct {
	EnvironmentName string `json:"environmentName"`
}

type NewDeploymentResponse struct {
	DeploymentId string `json:"deploymentId"`
	SessionToken string `json:"sessionToken"`
	UploadUrl    string `json:"uploadUrl"`
}

var (
	env   string
	stack string

	deployCmd = &cobra.Command{
		Use:   "deploy",
		Short: "Deploy a stack file",
		Long:  `Deploy a stack file to the specified environment.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Make sure that stack file path from flag actually exists in user directory
			if _, err := os.Stat(stack); err != nil {
				return err
			}

			newDeployment, err := CreateBrainDeployment()
			if err != nil {
				return err
			}

			fmt.Println("Deployment to brain created.")

			err = UploadStackFile(newDeployment.UploadUrl)
			if err != nil {
				return err
			}

			fmt.Println("Stack file uploaded!")

			return nil
		},
	}
)

func CreateBrainDeployment() (*NewDeploymentResponse, error) {
	deploymentReq := new(NewDeploymentRequest)
	deploymentReq.EnvironmentName = env

	body, err := json.Marshal(deploymentReq)
	if err != nil {
		return nil, err
	}

	url := viper.GetString("BrainUrl") + "/api/cli/deployment"
	res, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	responseData, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var newDeployment NewDeploymentResponse
	json.Unmarshal(responseData, &newDeployment)

	return &newDeployment, nil
}

func UploadStackFile(uploadUrl string) error {
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

func init() {
	// Do we **only** want to allow `dev`, `test`, `stage` & `prod` values here, or can they be named anything?
	deployCmd.PersistentFlags().StringVarP(&env, "env", "e", "test", "The environment to deploy to")
	deployCmd.MarkPersistentFlagRequired("env")

	deployCmd.PersistentFlags().StringVarP(&stack, "file", "f", "stack.yml", "The location of your stack file")
	deployCmd.MarkPersistentFlagRequired("file")
}
