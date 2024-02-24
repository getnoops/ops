//go:generate go run github.com/Khan/genqlient genqlient.yaml
package queries

import (
	"context"
	"strings"

	"github.com/Khan/genqlient/graphql"
	"github.com/getnoops/ops/pkg/config"
	"github.com/google/uuid"
)

type Queries interface {
	GetMemberOrganisations(ctx context.Context, page int, pageSize int) (*GetMemberOrganisationsMemberOrganisationsPagedOrganisationsOutput, error)
	GetCurrentOrganisation(ctx context.Context) (*Organisation, error)
	GetEnvironments(ctx context.Context, organisationId uuid.UUID, codes []string, page int, pageSize int) (*GetEnvironmentsEnvironmentsPagedEnvironmentsOutput, error)

	GetConfigs(ctx context.Context, organisationId uuid.UUID, class ConfigClass, page int, pageSize int) (*GetConfigsConfigsPagedConfigsOutput, error)
	GetConfig(ctx context.Context, organisationId uuid.UUID, code string) (*Config, error)
	UpdateConfig(ctx context.Context, input *UpdateConfigInput) (*uuid.UUID, error)

	GetContainerRepository(ctx context.Context, organisationId uuid.UUID, id uuid.UUID) (*ContainerRepository, error)
	CreateContainerRepository(ctx context.Context, organisationId uuid.UUID, id uuid.UUID, configId uuid.UUID, code string) (*uuid.UUID, error)
	DeleteContainerRepository(ctx context.Context, organisationId uuid.UUID, id uuid.UUID) (*uuid.UUID, error)
	LoginContainerRepository(ctx context.Context, organisationId uuid.UUID) (*AuthContainerRepository, error)

	GetApiKeys(ctx context.Context, organisationId uuid.UUID, page int, pageSize int) (*GetApiKeysApiKeysPagedApiKeysOutput, error)
	CreateApiKey(ctx context.Context, organisationId uuid.UUID) (*IdWithToken, error)
	UpdateApiKey(ctx context.Context, id uuid.UUID) (*IdWithToken, error)
	DeleteApiKey(ctx context.Context, id uuid.UUID) (*uuid.UUID, error)

	NewDeployment(ctx context.Context, organisationId uuid.UUID, environmentId uuid.UUID, configId uuid.UUID, configRevisionId uuid.UUID, revisionId uuid.UUID) (*uuid.UUID, error)
	GetDeploymentRevision(ctx context.Context, organisationId uuid.UUID, deploymentRevisionId uuid.UUID) (*DeploymentRevision, error)
}

type queries struct {
	organisationCode string
	client           graphql.Client
}

func (q *queries) GetMemberOrganisations(ctx context.Context, page int, pageSize int) (*GetMemberOrganisationsMemberOrganisationsPagedOrganisationsOutput, error) {
	resp, err := GetMemberOrganisations(ctx, q.client, page, pageSize)
	if err != nil {
		return nil, err
	}
	return resp.MemberOrganisations, nil
}

func (q *queries) GetCurrentOrganisation(ctx context.Context) (*Organisation, error) {
	orgs, err := q.GetMemberOrganisations(ctx, 1, 999)
	if err != nil {
		return nil, err
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

func (q *queries) GetEnvironments(ctx context.Context, organisationId uuid.UUID, codes []string, page int, pageSize int) (*GetEnvironmentsEnvironmentsPagedEnvironmentsOutput, error) {
	resp, err := GetEnvironments(ctx, q.client, organisationId, codes, page, pageSize)
	if err != nil {
		return nil, err
	}
	return resp.Environments, nil
}

func (q *queries) GetConfigs(ctx context.Context, organisationId uuid.UUID, class ConfigClass, page int, pageSize int) (*GetConfigsConfigsPagedConfigsOutput, error) {
	classes := []ConfigClass{class}

	resp, err := GetConfigs(ctx, q.client, organisationId, classes, page, pageSize)
	if err != nil {
		return nil, err
	}
	return resp.Configs, nil
}

func (q *queries) GetConfig(ctx context.Context, organisationId uuid.UUID, code string) (*Config, error) {
	resp, err := GetConfig(ctx, q.client, organisationId, code)
	if err != nil {
		return nil, err
	}
	return resp.Config, nil
}

func (q *queries) UpdateConfig(ctx context.Context, input *UpdateConfigInput) (*uuid.UUID, error) {
	resp, err := UpdateConfig(ctx, q.client, input)
	if err != nil {
		return nil, err
	}
	return &resp.UpdateConfig, nil
}

func (q *queries) GetContainerRepository(ctx context.Context, organisationId uuid.UUID, id uuid.UUID) (*ContainerRepository, error) {
	resp, err := GetContainerRepository(ctx, q.client, organisationId, id)
	if err != nil {
		return nil, err
	}
	return resp.ContainerRepository, nil
}

func (q *queries) CreateContainerRepository(ctx context.Context, organisationId uuid.UUID, id uuid.UUID, configId uuid.UUID, code string) (*uuid.UUID, error) {
	resp, err := CreateContainerRepository(ctx, q.client, organisationId, id, configId, code)
	if err != nil {
		return nil, err
	}
	return &resp.CreateContainerRepository, nil
}

func (q *queries) DeleteContainerRepository(ctx context.Context, organisationId uuid.UUID, id uuid.UUID) (*uuid.UUID, error) {
	resp, err := DeleteContainerRepository(ctx, q.client, organisationId, id)
	if err != nil {
		return nil, err
	}
	return &resp.DeleteContainerRepository, nil
}

func (q *queries) LoginContainerRepository(ctx context.Context, organisationId uuid.UUID) (*AuthContainerRepository, error) {
	resp, err := LoginContainerRepository(ctx, q.client, organisationId)
	if err != nil {
		return nil, err
	}
	return resp.LoginContainerRepository, nil
}

func (q *queries) GetApiKeys(ctx context.Context, organisationId uuid.UUID, page int, pageSize int) (*GetApiKeysApiKeysPagedApiKeysOutput, error) {
	resp, err := GetApiKeys(ctx, q.client, organisationId, page, pageSize)
	if err != nil {
		return nil, err
	}
	return resp.ApiKeys, nil
}

func (q *queries) CreateApiKey(ctx context.Context, organisationId uuid.UUID) (*IdWithToken, error) {
	id := uuid.New()
	resp, err := CreateApiKey(ctx, q.client, id, organisationId)
	if err != nil {
		return nil, err
	}
	return resp.CreateApiKey, nil
}

func (q *queries) UpdateApiKey(ctx context.Context, id uuid.UUID) (*IdWithToken, error) {
	resp, err := UpdateApiKey(ctx, q.client, id)
	if err != nil {
		return nil, err
	}
	return resp.UpdateApiKey, nil
}

func (q *queries) DeleteApiKey(ctx context.Context, id uuid.UUID) (*uuid.UUID, error) {
	resp, err := DeleteApiKey(ctx, q.client, id)
	if err != nil {
		return nil, err
	}
	return &resp.DeleteApiKey, nil
}

func (q *queries) NewDeployment(ctx context.Context, organisationId uuid.UUID, environmentId uuid.UUID, configId uuid.UUID, configRevisionId uuid.UUID, revisionId uuid.UUID) (*uuid.UUID, error) {
	id := uuid.New()
	resp, err := NewDeployment(ctx, q.client, organisationId, id, environmentId, configId, configRevisionId, revisionId)
	if err != nil {
		return nil, err
	}
	return &resp.NewDeployment, nil
}

func (q *queries) GetDeploymentRevision(ctx context.Context, organisationId uuid.UUID, deploymentRevisionId uuid.UUID) (*DeploymentRevision, error) {
	resp, err := GetDeploymentRevision(ctx, q.client, organisationId, deploymentRevisionId)
	if err != nil {
		return nil, err
	}
	return resp.DeploymentRevision, nil
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
