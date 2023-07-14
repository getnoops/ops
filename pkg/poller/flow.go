package poller

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/getnoops/ops/pkg/brain"
	"github.com/getnoops/ops/pkg/util"
)

// Specifies parameters to poll The Brain with until completion.
type WaitOptions struct {
	// The deployment to poll for
	DeploymentId string

	// Also known as `sessionToken`
	ExecToken *string

	// Client to make requests to brain
	BrainClient *brain.ClientWithResponses

	// The poller to use
	newPoller pollerFactory
}

// Polls The Brain until the session times out or the deployment completes.
func Wait(ctx context.Context, opts WaitOptions) error {
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
			return err
		}

		body := brain.CliPollRequest{CommandId: commandId}
		if opts.ExecToken != nil {
			body.ExecToken = opts.ExecToken
		}

		r, err := util.MakeBodyReaderFromType(body)
		if err != nil {
			return err
		}

		res, err := opts.BrainClient.PollForCommandsWithBodyWithResponse(ctx, opts.DeploymentId, "application/json", r)
		if err != nil {
			return err
		}

		var pollResponse brain.CliPollResponse
		json.Unmarshal(res.Body, &pollResponse)

		if err != nil {
			return err
		}

		if len(pollResponse.Commands) > 0 {
			commandType := "Previous"
			if commandId != nil {
				commandType = "New"
			}

			fmt.Printf("\n-----------------------\n")
			fmt.Printf("\n%s commands received: \n", commandType)
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
				return nil
			}
		}

		fmt.Println("\nWaiting for new commands...")
	}
}
