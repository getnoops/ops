package poller

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/getnoops/ops/pkg/brain"
	"github.com/spf13/viper"
)

// Specifies parameters to poll The Brain with until completion.
type WaitOptions struct {
	// The deployment to poll for
	deploymentId string

	// The execution command
	commandId string

	// Also known as `sessionToken`
	execToken string

	// The poller to use
	newPoller pollerFactory
}

// Polls The Brain until the session times out or the deployment completes.
func Wait(ctx context.Context, opts WaitOptions) (*[]brain.PollerQueueEntry, error) {
	seconds := 10
	checkInterval := time.Duration(seconds) * time.Second

	minutes := 30
	expiresIn := time.Duration(minutes) * time.Minute

	makePoller := opts.newPoller
	if makePoller == nil {
		makePoller = newPoller
	}
	_, poll := makePoller(ctx, checkInterval, expiresIn)

	for {
		if err := poll.Wait(); err != nil {
			return nil, err
		}

		body := brain.CliPollRequest{CommandId: &opts.commandId, ExecToken: &opts.execToken}

		url := viper.GetString("BrainUrl")
		req, err := brain.NewPollForCommandsRequest(url, opts.deploymentId, body)

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

		var pollResponse brain.CliPollResponse
		json.Unmarshal(resData, &pollResponse)

		if err != nil {
			return nil, err
		}

		return &pollResponse.Commands, nil
	}
}
