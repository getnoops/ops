package poller

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/getnoops/ops/pkg/brain"
)

func makeRequestToPollEndpoint(ctx context.Context, opts WaitOptions) (*brain.CliPollResponse, *brain.PollForCommandsResponse, error) {
	body := brain.CliPollRequest{CommandId: commandId}
	if opts.ExecToken != nil {
		body.ExecToken = opts.ExecToken
	}

	res, err := brain.Client.PollForCommandsWithResponse(ctx, opts.DeploymentId, body)
	if err != nil {
		return nil, res, err
	}

	var pollResponse brain.CliPollResponse
	json.Unmarshal(res.Body, &pollResponse)
	if err != nil {
		return nil, res, err
	}

	return &pollResponse, res, nil
}

func makeRequestToDockerLoginEndpoint(ctx context.Context, deploymentId, artifactId string) (*brain.DockerLoginResponse, error) {
	fmt.Println("\nGetting docker login credentials")

	res, err := brain.Client.GetDockerLoginWithResponse(ctx, deploymentId, artifactId)
	if err != nil {
		return nil, err
	}

	var dockerLogin brain.DockerLoginResponse
	err = json.Unmarshal(res.Body, &dockerLogin)
	if err != nil {
		return nil, err
	}

	return &dockerLogin, nil
}
