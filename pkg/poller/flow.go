package poller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/getnoops/ops/pkg/brain"
	"github.com/spf13/viper"
)

// Specifies parameters to poll The Brain with until completion.
type WaitOptions struct {
	// The deployment to poll for
	DeploymentId string

	// Also known as `sessionToken`
	ExecToken string

	// The poller to use
	newPoller pollerFactory
}

// Polls The Brain until the session times out or the deployment completes.
func Wait(ctx context.Context, opts WaitOptions) (*brain.PollerQueueEntry, error) {
	var commandId *string

	seconds := 10
	checkInterval := time.Duration(seconds) * time.Second

	minutes := 60
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

		body := brain.CliPollRequest{CommandId: commandId, ExecToken: &opts.ExecToken}

		url := viper.GetString("BrainUrl")
		req, err := brain.NewPollForCommandsRequest(url, opts.DeploymentId, body)

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

		if len(pollResponse.Commands) > 0 {
			fmt.Printf("\n-----------------------\n")
			fmt.Println("\nNew commands received:")
			fmt.Printf("\n-----------------------\n")

			for _, c := range pollResponse.Commands {
				fmt.Printf("\nCommand order: %d", c.SeqOrder)
				fmt.Printf("\nCommand type: %s", c.CmdType)
				fmt.Printf("\nCommand: %s", c.Command)
				fmt.Printf("\n\n-----------------------\n")
			}

			lastCommand := pollResponse.Commands[len(pollResponse.Commands)-1]
			commandId = lastCommand.Id

			if lastCommand.CmdType == brain.DEPLOYMENTFINISHED {
				return &lastCommand, nil
			}
		}

		fmt.Println("\nWaiting for new commands...")
	}
}
