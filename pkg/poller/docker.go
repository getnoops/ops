package poller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/avast/retry-go/v4"
	dockerTypes "github.com/docker/docker/api/types"
	docker "github.com/docker/docker/client"
	"github.com/getnoops/ops/pkg/brain"
)

type PushDockerImageCommandInfo struct {
	ArtifactId string `json:"artifactId"`

	Img string `json:"img"`

	Tag string `json:"tag"`

	DeploymentId string `json:"deploymentId"`

	Type brain.PollerQueueEntryCmdType `json:"type"`
}

func formatDockerCommandInfo(commandMsg string) (*PushDockerImageCommandInfo, error) {
	var dockerCommandInfo PushDockerImageCommandInfo

	err := json.Unmarshal([]byte(commandMsg), &dockerCommandInfo)
	if err != nil {
		return nil, err
	}

	return &dockerCommandInfo, nil
}

// func tagDockerImageWithEcrUrl(dockerCommandInfo *PushDockerImageCommandInfo, dockerLogin *brain.DockerLoginResponse) error {
// 	userProvidedImage := fmt.Sprintf("%s:%s", dockerCommandInfo.Img, dockerCommandInfo.Tag)
// 	fmt.Printf("Tagging image [%s] with [%s]", userProvidedImage, dockerLogin.Url)

// 	cmd := exec.Command("docker", "tag", userProvidedImage, dockerLogin.Url)
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr
// 	err := cmd.Run()
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func pushImage(ctx context.Context, dockerClient *docker.Client, docker *brain.DockerLoginResponse) error {
	closer, err := dockerClient.ImagePush(ctx, docker.Url, dockerTypes.ImagePushOptions{RegistryAuth: docker.UserName + ":" + docker.Password})
	if err != nil {
		return err
	}
	closer.Close()

	return nil
}

func pushImageWithRetry(ctx context.Context, b brain.Manager, dockerClient *docker.Client, deploymentId string, docker *brain.DockerLoginResponse) error {
	err := retry.Do(
		func() error {
			err := pushImage(ctx, dockerClient, docker)
			return err
		},
		retry.Attempts(3),
		retry.OnRetry(func(n uint, err error) {
			log.Printf("Unable to push docker image to ECR. Retrying request after error: %v", err)
		}),
	)
	if err != nil {
		e := err.Error()
		notifyDockerUploadBody := brain.NotifyUploadCompleteRequest{Success: false, Error: &e}
		b.NotifyDockerUploadCompleted(ctx, deploymentId, notifyDockerUploadBody)
		return err
	}

	return nil
}

func pushDockerImageToECR(ctx context.Context, b brain.Manager, command *brain.PollerQueueEntry, deploymentId string) error {
	fmt.Println("\nStarting process to push your docker image to ECR...")

	dockerCli, err := docker.NewClientWithOpts(docker.FromEnv, docker.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	dockerCommandInfo, err := formatDockerCommandInfo(command.Command)
	if err != nil {
		return err
	}

	dockerLogin, err := b.GetECRCredentials(ctx, deploymentId, dockerCommandInfo.ArtifactId)
	if err != nil {
		return err
	}

	// TODO: Build docker image

	// registryUrl := fmt.Sprintf("https://%s", dockerLogin.Url)

	userImage := dockerCommandInfo.Img + ":" + dockerCommandInfo.Tag
	fmt.Printf("Tagging image [%s] with [%s]", userImage, dockerLogin.Url)
	err = dockerCli.ImageTag(ctx, userImage, dockerLogin.Url)
	if err != nil {
		return err
	}

	// err = tagDockerImageWithEcrUrl(dockerCommandInfo, dockerLogin)
	// if err != nil {
	// 	return err
	// }

	// err = loginToDocker(dockerLogin)
	// if err != nil {
	// 	return err
	// }

	err = pushImageWithRetry(ctx, b, dockerCli, deploymentId, dockerLogin)

	if err != nil {
		return err
	}

	fmt.Println("\nSuccessfully pushed your image to ECR!")

	notifyUploadCompleteBody := brain.NotifyUploadCompleteRequest{Success: true}
	err = b.NotifyDockerUploadCompleted(ctx, deploymentId, notifyUploadCompleteBody)
	if err != nil {
		return err
	}

	fmt.Println("\nBrain notified that Docker Image has been uploaded.")

	return nil
}
