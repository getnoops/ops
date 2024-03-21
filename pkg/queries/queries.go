//go:generate go run github.com/Khan/genqlient genqlient.yaml
package queries

import (
	"context"
	"fmt"
	"strings"

	"github.com/Khan/genqlient/graphql"
	"github.com/getnoops/ops/pkg/config"
	"github.com/google/uuid"
)

type Queries interface {
	GetMemberOrganisations(ctx context.Context, page int, pageSize int) (*GetMemberOrganisationsMemberOrganisationsPagedOrganisationsOutput, error)
	GetCurrentOrganisation(ctx context.Context) (*Organisation, error)
	GetEnvironments(ctx context.Context, organisationId uuid.UUID, codes []string, states []StackState, page int, pageSize int) (*GetEnvironmentsEnvironmentsPagedEnvironmentsOutput, error)

	GetConfigs(ctx context.Context, organisationId uuid.UUID, classes []ConfigClass, page int, pageSize int) (*GetConfigsConfigsPagedConfigsOutput, error)
	GetConfig(ctx context.Context, organisationId uuid.UUID, code string) (*Config, error)
	CreateConfig(ctx context.Context, organisationId uuid.UUID, id uuid.UUID, name string, code string, class ConfigClass) (*uuid.UUID, error)
	UpdateConfig(ctx context.Context, input *UpdateConfigInput) (*uuid.UUID, error)

	CreateContainerRepository(ctx context.Context, organisationId uuid.UUID, id uuid.UUID, configId uuid.UUID, code string) (*uuid.UUID, error)
	DeleteContainerRepository(ctx context.Context, organisationId uuid.UUID, id uuid.UUID) (*uuid.UUID, error)
	LoginContainerRepository(ctx context.Context, organisationId uuid.UUID) (*AuthContainerRepository, error)

	CreateSecret(ctx context.Context, organisationId uuid.UUID, id uuid.UUID, configId uuid.UUID, environmentId uuid.UUID, code string, value string) (*uuid.UUID, error)
	DeleteSecret(ctx context.Context, organisationId uuid.UUID, id uuid.UUID) (*uuid.UUID, error)

	GetApiKeys(ctx context.Context, organisationId uuid.UUID, page int, pageSize int) (*GetApiKeysApiKeysPagedApiKeysOutput, error)
	CreateApiKey(ctx context.Context, organisationId uuid.UUID) (*IdWithToken, error)
	UpdateApiKey(ctx context.Context, id uuid.UUID) (*IdWithToken, error)
	DeleteApiKey(ctx context.Context, id uuid.UUID) (*uuid.UUID, error)

	NewDeployment(ctx context.Context, organisationId uuid.UUID, deploymentId uuid.UUID, environmentId uuid.UUID, configId uuid.UUID, configRevisionId uuid.UUID, revisionId uuid.UUID) (*uuid.UUID, error)
	DeleteDeployment(ctx context.Context, organisationId uuid.UUID, deploymentId uuid.UUID) (*uuid.UUID, error)
	GetDeploymentRevision(ctx context.Context, organisationId uuid.UUID, deploymentRevisionId uuid.UUID) (*DeploymentRevision, error)
	GetDeployment(ctx context.Context, organisationId uuid.UUID, deploymentId uuid.UUID) (*Deployment, error)
}

type queries struct {
	organisationCode string
	client           graphql.Client
}

func (q *queries) GetMemberOrganisations(ctx context.Context, page int, pageSize int) (*GetMemberOrganisationsMemberOrganisationsPagedOrganisationsOutput, error) {
	resp, err := GetMemberOrganisations(ctx, q.client, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("GetMemberOrganisations unexpected response: %v", err)
	}
	return resp.MemberOrganisations, nil
}

func (q *queries) GetCurrentOrganisation(ctx context.Context) (*Organisation, error) {
	orgs, err := q.GetMemberOrganisations(ctx, 1, 999)
	if err != nil {
		return nil, fmt.Errorf("GetCurrentOrganisation unexpected response: %v", err)

	}

	if len(q.organisationCode) == 0 && len(orgs.Items) == 1 {
		return orgs.Items[0], nil
	}

	for _, org := range orgs.Items {
		if strings.EqualFold(org.Code, q.organisationCode) {
			return org, nil
		}
	}

	return nil, config.ErrNoOrganisation
}

func (q *queries) GetEnvironments(ctx context.Context, organisationId uuid.UUID, codes []string, states []StackState, page int, pageSize int) (*GetEnvironmentsEnvironmentsPagedEnvironmentsOutput, error) {
	resp, err := GetEnvironments(ctx, q.client, organisationId, codes, states, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("GetEnvironments unexpected response: %v", err)

	}
	return resp.Environments, nil
}

func (q *queries) GetConfigs(ctx context.Context, organisationId uuid.UUID, classes []ConfigClass, page int, pageSize int) (*GetConfigsConfigsPagedConfigsOutput, error) {
	resp, err := GetConfigs(ctx, q.client, organisationId, classes, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("GetConfigs unexpected response: %v", err)

	}
	return resp.Configs, nil
}

func (q *queries) GetConfig(ctx context.Context, organisationId uuid.UUID, code string) (*Config, error) {
	resp, err := GetConfig(ctx, q.client, organisationId, code)
	if err != nil {
		return nil, fmt.Errorf("GetConfig unexpected response: %v", err)

	}
	return resp.Config, nil
}

func (q *queries) CreateConfig(ctx context.Context, organisationId uuid.UUID, id uuid.UUID, name string, code string, class ConfigClass) (*uuid.UUID, error) {
	resp, err := CreateConfig(ctx, q.client, organisationId, id, code, class, name)
	if err != nil {
		return nil, fmt.Errorf("CreateConfig unexpected response: %v", err)

	}
	return &resp.CreateConfig, nil

}

func (q *queries) UpdateConfig(ctx context.Context, input *UpdateConfigInput) (*uuid.UUID, error) {
	resp, err := UpdateConfig(ctx, q.client, input)
	if err != nil {
		return nil, fmt.Errorf("UpdateConfig unexpected response: %v", err)

	}
	return &resp.UpdateConfig, nil
}

func (q *queries) CreateContainerRepository(ctx context.Context, organisationId uuid.UUID, id uuid.UUID, configId uuid.UUID, code string) (*uuid.UUID, error) {
	resp, err := CreateContainerRepository(ctx, q.client, organisationId, id, configId, code)
	if err != nil {
		return nil, fmt.Errorf("CreateContainerRepository unexpected response: %v", err)

	}
	return &resp.CreateContainerRepository, nil
}

func (q *queries) DeleteContainerRepository(ctx context.Context, organisationId uuid.UUID, id uuid.UUID) (*uuid.UUID, error) {
	resp, err := DeleteContainerRepository(ctx, q.client, organisationId, id)
	if err != nil {
		return nil, fmt.Errorf("DeleteContainerRepository unexpected response: %v", err)

	}
	return &resp.DeleteContainerRepository, nil
}

func (q *queries) LoginContainerRepository(ctx context.Context, organisationId uuid.UUID) (*AuthContainerRepository, error) {
	resp, err := LoginContainerRepository(ctx, q.client, organisationId)
	if err != nil {
		return nil, fmt.Errorf("LoginContainerRepository unexpected response: %v", err)

	}
	return resp.LoginContainerRepository, nil
}

func (q *queries) CreateSecret(ctx context.Context, organisationId uuid.UUID, id uuid.UUID, configId uuid.UUID, environmentId uuid.UUID, code string, value string) (*uuid.UUID, error) {
	resp, err := CreateSecret(ctx, q.client, organisationId, id, configId, environmentId, code, value)
	if err != nil {
		return nil, fmt.Errorf("CreateSecret unexpected response: %v", err)

	}
	return &resp.CreateSecret, nil
}

func (q *queries) DeleteSecret(ctx context.Context, organisationId uuid.UUID, id uuid.UUID) (*uuid.UUID, error) {
	resp, err := DeleteSecret(ctx, q.client, organisationId, id)
	if err != nil {
		return nil, fmt.Errorf("DeleteSecret unexpected response: %v", err)

	}
	return &resp.DeleteSecret, nil
}

func (q *queries) GetApiKeys(ctx context.Context, organisationId uuid.UUID, page int, pageSize int) (*GetApiKeysApiKeysPagedApiKeysOutput, error) {
	resp, err := GetApiKeys(ctx, q.client, organisationId, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("GetApiKeys unexpected response: %v", err)

	}
	return resp.ApiKeys, nil
}

func (q *queries) CreateApiKey(ctx context.Context, organisationId uuid.UUID) (*IdWithToken, error) {
	id := uuid.New()
	resp, err := CreateApiKey(ctx, q.client, id, organisationId)
	if err != nil {
		return nil, fmt.Errorf("CreateApiKey unexpected response: %v", err)

	}
	return resp.CreateApiKey, nil
}

func (q *queries) UpdateApiKey(ctx context.Context, id uuid.UUID) (*IdWithToken, error) {
	resp, err := UpdateApiKey(ctx, q.client, id)
	if err != nil {
		return nil, fmt.Errorf("UpdateApiKey unexpected response: %v", err)

	}
	return resp.UpdateApiKey, nil
}

func (q *queries) DeleteApiKey(ctx context.Context, id uuid.UUID) (*uuid.UUID, error) {
	resp, err := DeleteApiKey(ctx, q.client, id)
	if err != nil {
		return nil, fmt.Errorf("DeleteApiKey unexpected response: %v", err)

	}
	return &resp.DeleteApiKey, nil
}

func (q *queries) NewDeployment(ctx context.Context, organisationId uuid.UUID, deploymentId uuid.UUID, environmentId uuid.UUID, configId uuid.UUID, configRevisionId uuid.UUID, revisionId uuid.UUID) (*uuid.UUID, error) {
	resp, err := NewDeployment(ctx, q.client, organisationId, deploymentId, environmentId, configId, configRevisionId, revisionId)
	if err != nil {
		return nil, fmt.Errorf("NewDeployment unexpected response: %v", err)

	}
	return &resp.NewDeployment, nil
}

func (q *queries) DeleteDeployment(ctx context.Context, organisationId uuid.UUID, deploymentId uuid.UUID) (*uuid.UUID, error) {
	resp, err := DeleteDeployment(ctx, q.client, organisationId, deploymentId)
	if err != nil {
		return nil, fmt.Errorf("DeleteDeployment unexpected response: %v", err)

	}
	return &resp.DeleteDeployment, nil
}

func (q *queries) GetDeploymentRevision(ctx context.Context, organisationId uuid.UUID, deploymentRevisionId uuid.UUID) (*DeploymentRevision, error) {
	resp, err := GetDeploymentRevision(ctx, q.client, organisationId, deploymentRevisionId)
	if err != nil {
		return nil, fmt.Errorf("GetDeploymentRevision unexpected response: %v", err)

	}
	return resp.DeploymentRevision, nil
}

func (q *queries) GetDeployment(ctx context.Context, organisationId uuid.UUID, deploymentId uuid.UUID) (*Deployment, error) {
	resp, err := GetDeployment(ctx, q.client, organisationId, deploymentId)
	if err != nil {
		return nil, fmt.Errorf("GetDeployment unexpected response: %v", err)

	}
	return resp.Deployment, nil
}

func New[C any, T any](ctx context.Context, cfg *config.NoOps[C, T]) (Queries, error) {
	httpClient, err := cfg.NewHttpClient(ctx)
	if err != nil {
		return nil, err
	}

	organisationCode := cfg.GetOrganisationCode()

	client := graphql.NewClient(cfg.Api.GraphQL, httpClient)
	return &queries{
		organisationCode: organisationCode,
		client:           client,
	}, nil
}
