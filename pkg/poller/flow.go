package poller

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
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
	var firstPass = true

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
		if !firstPass {
			if err := poll.Wait(); err != nil {
				return err
			}
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

				if c.CmdType == brain.PUSHDOCKERIMAGE || c.CmdType == brain.UPLOADSTATICFILE {
					handleWorkCommand(&c, opts.DeploymentId, opts.BrainClient, ctx)
				}

			}

			lastCommand := pollResponse.Commands[len(pollResponse.Commands)-1]
			commandId = lastCommand.Id

			if lastCommand.CmdType == brain.DEPLOYMENTFINISHED {
				return nil
			}
		}

		fmt.Println("\nWaiting for events...")

		if firstPass {
			firstPass = false
		}
	}
}

type PushDockerImageCommandInfo struct {
	ArtifactId string `json:"artifactId"`

	Img string `json:"img"`

	Tag string `json:"tag"`

	DeploymentId string `json:"deploymentId"`

	Type brain.PollerQueueEntryCmdType `json:"type"`
}

func handleWorkCommand(command *brain.PollerQueueEntry, deploymentId string, client *brain.ClientWithResponses, ctx context.Context) error {
	if command.CmdType == brain.PUSHDOCKERIMAGE {
		fmt.Println("\nStarting process to push your docker image to ECR...")

		var dockerCommandInfo PushDockerImageCommandInfo
		err := json.Unmarshal([]byte(command.Command), &dockerCommandInfo)
		if err != nil {
			return err
		}

		fmt.Println("\nGetting docker login credentials")
		res, err := client.GetDockerLoginWithResponse(ctx, deploymentId, dockerCommandInfo.ArtifactId)
		if err != nil {
			return err
		}

		var dockerLogin brain.DockerLoginResponse
		err = json.Unmarshal(res.Body, &dockerLogin)
		if err != nil {
			return err
		}

		fmt.Println(dockerLogin.Url)

		userProvidedImage := fmt.Sprintf("%s:%s", dockerCommandInfo.Img, dockerCommandInfo.Tag)

		fmt.Printf("Tagging image [%s] with [%s] \n", userProvidedImage, dockerLogin.Url)
		cmd := exec.Command("docker", "tag", userProvidedImage, dockerLogin.Url)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			return err
		}

		// TODO: Investigate if there's a better way to add the password here.
		// ``WARNING! Using --password via the CLI is insecure. Use --password-stdin
		fmt.Println("\nLogging in to docker")
		registryUrl := fmt.Sprintf("https://%s", dockerLogin.Url)
		cmd = exec.Command("docker", "login", "--username", dockerLogin.UserName, "--password", dockerLogin.Password, registryUrl)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			return err
		}

		cmd = exec.Command("docker", "push", dockerLogin.Url)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			return err
		}

		fmt.Println("\nSuccessfully pushed your image to ECR!")
		fmt.Printf("\n-----------------------\n")
	} else if command.CmdType == brain.UPLOADSTATICFILE {
		// Can be implemented later, we're not handling static files for MVP
		context.TODO()
	}

	return nil
}
