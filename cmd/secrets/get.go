package secrets

import (
	"context"
	"fmt"

	"github.com/getnoops/ops/pkg/config"
	"github.com/getnoops/ops/pkg/queries"
	"github.com/getnoops/ops/pkg/util"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Secret struct {
	Id          uuid.UUID         `json:"id"`
	Code        string            `json:"code"`
	State       string            `json:"state"`
	Environment string            `json:"environment"`
	Outputs     map[string]string `json:"outputs"`
}

type GetConfig struct {
}

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "get [compute] [env] [code]",
		Short:  "Will get a container repository for a given compute repository",
		Args:   cobra.ExactArgs(3),
		PreRun: util.BindPreRun,
		RunE: func(cmd *cobra.Command, args []string) error {
			configCode := args[0]
			envCode := args[1]
			code := args[2]

			ctx := cmd.Context()
			return Get(ctx, configCode, envCode, code)
		},
		ValidArgs: []string{"compute", "code"},
	}
	return cmd
}

func Get(ctx context.Context, computeCode string, environmentCode string, code string) error {
	cfg, err := config.New[GetConfig, *Secret](ctx, viper.GetViper())
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
		cfg.WriteStderr("failed to get configs")
		return nil
	}

	secret, err := GetSecret(config.Secrets, environmentCode, code)
	if err != nil {
		cfg.WriteStderr("failed to get secret")
		return nil
	}

	outputs := map[string]string{}
	for _, output := range secret.Stack.Outputs {
		outputs[output.Output_key] = output.Output_value
	}

	cfg.WriteObject(&Secret{
		Id:          secret.Id,
		Code:        secret.Code,
		State:       string(secret.State),
		Environment: secret.Environment.Code,
		Outputs:     outputs,
	})
	return nil
}

func GetSecret(secrets []*queries.SecretItem, environmentCode string, code string) (*queries.SecretItem, error) {
	for _, secret := range secrets {
		if secret.Code == code && secret.Environment.Code == environmentCode {
			return secret, nil
		}
	}

	return nil, fmt.Errorf("secret not found")
}
