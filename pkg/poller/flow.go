package poller

import (
	"context"
	"fmt"
	"time"

	"github.com/getnoops/ops/pkg/brain"
)

type PollConfig struct {
	// The time to wait between each poll (Seconds)
	Interval int

	// The total time to run the poller until it automatically exits (Minutes)
	Expiry int
}

// Specifies parameters to poll The Brain with until completion.
type WaitOptions struct {
	// The deployment to poll for
	DeploymentId string

	// Also known as `sessionToken`
	ExecToken *string

	// Config for the poller
	PollerConfig PollConfig
}

var (
	// The last event command ID.
	// Used to get brain events **after** a specific event
	commandId *string

	// Ignores the `interval` value on the first pass of the poll.
	// Otherwise the first time the function is called, the poller will wait before sending a request
	firstPass = true
)

// Polls The Brain until the session times out or the deployment completes.
func Wait(ctx context.Context, opts WaitOptions) error {
	interval := formatIntToTime(opts.PollerConfig.Interval, time.Second)
	expiry := formatIntToTime(opts.PollerConfig.Expiry, time.Minute)

	_, poll := newPoller(ctx, interval, expiry)

	for {
		if !firstPass {
			if err := poll.Wait(); err != nil {
				return err
			}
		}

		pollResponse, httpResponse, err := makeRequestToPollEndpoint(ctx, opts)
		if httpResponse.StatusCode() == 409 {
			printLineBreak()
			fmt.Printf("\nDeployment has been completed.\n")
			printLineBreak()
			return nil
		}
		if err != nil {
			return err
		}

		if len(pollResponse.Commands) > 0 {
			for _, c := range pollResponse.Commands {
				printLineBreak()
				fmt.Printf("\nCommand order: %d", c.SeqOrder)
				fmt.Printf("\nCommand type: %s", c.CmdType)
				fmt.Printf("\nCommand: %s\n", c.Command)

				if c.CmdType == brain.PUSHDOCKERIMAGE {
					pushDockerImageToECR(ctx, &c, opts.DeploymentId)
				} else if c.CmdType == brain.UPLOADSTATICFILE {
					// Can be implemented later, we're not handling static files for MVP
					context.TODO()
				}
			}

			lastCommand := pollResponse.Commands[len(pollResponse.Commands)-1]
			commandId = lastCommand.Id

			if lastCommand.CmdType == brain.DEPLOYMENTFINISHED {
				return nil
			}
		}

		printLineBreak()
		fmt.Printf("\nWaiting for events...\n")

		if firstPass {
			firstPass = false
		}
	}
}

func formatIntToTime(value int, unit time.Duration) time.Duration {
	return time.Duration(value) * unit
}

func printLineBreak() {
	fmt.Printf("\n-----------------------\n")
}
