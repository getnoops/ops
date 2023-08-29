package brain

import (
	"context"
	"fmt"
	"net/http"
)

type Manager interface {
	CreateDeployment(ctx context.Context, envName string) (*CreateDeploymentResponse, error)
	NotifyStackFileUploadCompleted(ctx context.Context, deploymentId string, body NotifyUploadCompleteRequest) error
	NotifyDockerUploadCompleted(ctx context.Context, deploymentId string, body NotifyUploadCompleteRequest) error
	ListActiveDeployments(ctx context.Context) (*[]ActiveDeployment, error)
	GetECRCredentials(ctx context.Context, deploymentId, artifactId string) (*DockerLoginResponse, error)
	PollForCommands(ctx context.Context, deploymentId string, commandId, execToken *string) (*CliPollResponse, *int, error)
}

type manager struct {
	cli ClientWithResponsesInterface
}

func (m *manager) CreateDeployment(ctx context.Context, envName string) (*CreateDeploymentResponse, error) {
	rsp, err := m.cli.CreateNewDeploymentWithResponse(ctx, CreateDeploymentRequest{EnvironmentName: envName})
	if err != nil {
		return nil, err
	}

	if rsp.JSON200 != nil {
		return rsp.JSON200, nil
	}

	return nil, fmt.Errorf("CreateDeployment unexpected response: %s", rsp.Status())
}

func (m *manager) NotifyStackFileUploadCompleted(ctx context.Context, deploymentId string, body NotifyUploadCompleteRequest) error {
	rsp, err := m.cli.NotifyStackFileUploadCompletedWithResponse(ctx, deploymentId, body)
	if err != nil {
		return err
	}

	if rsp.StatusCode() != 200 {
		return fmt.Errorf("NotifyStackFileUploadCompleted unexpected response: %s", rsp.Status())
	}

	return nil
}

func (m *manager) NotifyDockerUploadCompleted(ctx context.Context, deploymentId string, body NotifyUploadCompleteRequest) error {
	rsp, err := m.cli.NotifyDockerUploadCompletedWithResponse(ctx, deploymentId, body)
	if err != nil {
		return err
	}

	if rsp.StatusCode() != 200 {
		return fmt.Errorf("NotifyDockerUploadCompleted unexpected response: %s", rsp.Status())
	}

	return nil
}

func (m *manager) ListActiveDeployments(ctx context.Context) (*[]ActiveDeployment, error) {
	rsp, err := m.cli.ListActiveDeploymentsWithResponse(ctx)
	if err != nil {
		return nil, err
	}

	if rsp.JSON200 != nil {
		return rsp.JSON200, nil
	}

	return nil, fmt.Errorf("ListActiveDeployments unexpected response: %s", rsp.Status())
}

func (m *manager) GetECRCredentials(ctx context.Context, deploymentId, artifactId string) (*DockerLoginResponse, error) {
	fmt.Println("\nGetting credentials to push image")

	rsp, err := m.cli.GetDockerLoginWithResponse(ctx, deploymentId, artifactId)
	if err != nil {
		return nil, err
	}

	if rsp.JSON200 != nil {
		return rsp.JSON200, nil
	}

	return nil, fmt.Errorf("GetECRCredentials unexpected response: %s", rsp.Status())
}

func (m *manager) PollForCommands(ctx context.Context, deploymentId string, commandId, execToken *string) (*CliPollResponse, *int, error) {
	body := CliPollRequest{CommandId: commandId}
	if execToken != nil {
		body.ExecToken = execToken
	}

	res, err := m.cli.PollForCommandsWithResponse(ctx, deploymentId, body)
	if err != nil {
		return nil, nil, err
	}

	code := res.StatusCode()

	if res.JSON200 != nil {
		return res.JSON200, &code, nil
	}

	if code == 409 {
		return nil, &code, nil
	}

	return nil, &code, fmt.Errorf("PollForCommands unexpected response: %v", code)
}

func NewManager(brainUrl string, httpClient *http.Client) (Manager, error) {
	cli, err := NewClientWithResponses(brainUrl, WithHTTPClient(httpClient))

	if err != nil {
		return nil, err
	}

	return &manager{
		cli: cli,
	}, nil
}
