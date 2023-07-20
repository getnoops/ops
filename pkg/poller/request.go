package poller

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/getnoops/ops/pkg/brain"
	"github.com/getnoops/ops/pkg/util"
)

func makeRequestToPollEndpoint(ctx context.Context, opts WaitOptions) (*brain.CliPollResponse, error) {
	body := brain.CliPollRequest{CommandId: commandId}
	if opts.ExecToken != nil {
		body.ExecToken = opts.ExecToken
	}

	r, err := util.MakeBodyReaderFromType(body)
	if err != nil {
		return nil, err
	}

	res, err := brain.Client.PollForCommandsWithBodyWithResponse(ctx, opts.DeploymentId, "application/json", r)
	if err != nil {
		return nil, err
	}

	var pollResponse brain.CliPollResponse
	json.Unmarshal(res.Body, &pollResponse)
	if err != nil {
		return nil, err
	}

	return &pollResponse, nil
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
