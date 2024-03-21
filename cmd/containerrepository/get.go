package containerrepository

import (
	"context"
	"fmt"

	"github.com/contextcloud/goutils/xmap"
	"github.com/getnoops/ops/pkg/config"
	"github.com/getnoops/ops/pkg/queries"
	"github.com/getnoops/ops/pkg/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Repository struct {
	RepositoryUri  string `json:"repository_uri"`
	RepositoryName string `json:"repository_name"`
	Username       string `json:"username"`
}

type GetConfig struct {
}

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "get [compute] [code]",
		Short:  "Will get a container repository for a given compute repository",
		Args:   cobra.ExactArgs(2),
		PreRun: util.BindPreRun,
		RunE: func(cmd *cobra.Command, args []string) error {
			configCode := args[0]
			code := args[1]

			ctx := cmd.Context()
			return Get(ctx, configCode, code)
		},
		ValidArgs: []string{"compute", "code"},
	}
	return cmd
}

func Get(ctx context.Context, computeCode string, code string) error {
	cfg, err := config.New[GetConfig, Repository](ctx, viper.GetViper())
	if err != nil {
		return err
	}

	q, err := queries.New(ctx, cfg)
	if err != nil {
		return err
	}

	organisation, err := q.GetCurrentOrganisation(ctx)
	if err == config.ErrNoOrganisation {
		cfg.WriteStderr("no organisation set")
		return nil
	}
	if err != nil {
		return err
	}

	config, err := q.GetConfig(ctx, organisation.Id, computeCode)
	if err != nil {
		return err
	}

	repository, err := GetRepository(config.ContainerRepositories, code)
	if err != nil {
		return err
	}

	outputs := map[string]string{}
	for _, output := range repository.Stack.Outputs {
		outputs[output.Output_key] = output.Output_value
	}

	result := Repository{
		Username: "AWS",
	}
	if err := xmap.Decode(outputs, &result); err != nil {
		cfg.WriteStderr(fmt.Sprintf("failed to decode stack outputs: {outputs: %v} {err: %v}", outputs, err))
		return nil
	}

	cfg.WriteObject(result)
	return nil
}

func GetRepository(repositories []*queries.ContainerRepositoryItem, code string) (*queries.ContainerRepositoryItem, error) {
	for _, repository := range repositories {
		if repository.Code == code {
			return repository, nil
		}
	}

	return nil, fmt.Errorf("repository not found")
}
