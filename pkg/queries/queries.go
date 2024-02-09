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
	GetCurrentOrganisation(ctx context.Context) (*Organisation, error)
	GetConfigs(ctx context.Context, organisationId uuid.UUID, class ConfigClass, page int, pageSize int) (*GetConfigsConfigsPagedConfigsOutput, error)
	GetConfig(ctx context.Context, organisationId uuid.UUID, code string) (*ConfigWithRevisions, error)
	GetEnvironments(ctx context.Context, organisationId uuid.UUID, codes []string, page int, pageSize int) (*GetEnvironmentsEnvironmentsPagedEnvironmentsOutput, error)
	GetMemberOrganisations(ctx context.Context, page int, pageSize int) (*GetMemberOrganisationsMemberOrganisationsPagedOrganisationsOutput, error)
	NewDeployment(ctx context.Context, organisationId uuid.UUID, environmentId uuid.UUID, configId uuid.UUID, configRevisionId uuid.UUID, revisionId uuid.UUID) (uuid.UUID, error)

	CreateContainerRepository(ctx context.Context, organisationId uuid.UUID, configId uuid.UUID, code string) (uuid.UUID, error)
	DeleteContainerRepository(ctx context.Context, organisationId uuid.UUID, id uuid.UUID) (uuid.UUID, error)
	LoginContainerRepository(ctx context.Context, organisationId uuid.UUID, id uuid.UUID) (*AuthContainerRepository, error)
}

type queries struct {
	organisation string
	userId       uuid.UUID
	client       graphql.Client
}

func (q *queries) GetCurrentOrganisation(ctx context.Context) (*Organisation, error) {
	if q.organisation == "" {
		return nil, config.ErrNoOrganisation
	}

	orgs, err := q.GetMemberOrganisations(ctx, 1, 999)
	if err != nil {
		return nil, err
	}

	for _, org := range orgs.Items {
		if strings.EqualFold(org.Name, q.organisation) {
			return &org, nil
		}
	}

	return nil, config.ErrNoOrganisation
}

func (q *queries) GetConfigs(ctx context.Context, organisationId uuid.UUID, class ConfigClass, page int, pageSize int) (*GetConfigsConfigsPagedConfigsOutput, error) {
	resp, err := GetConfigs(ctx, q.client, organisationId, class, page, pageSize)
	if err != nil {
		return nil, err
	}
	return &resp.Configs, nil
}

func (q *queries) GetConfig(ctx context.Context, organisationId uuid.UUID, code string) (*ConfigWithRevisions, error) {
	resp, err := GetConfig(ctx, q.client, organisationId, code)
	if err != nil {
		return nil, err
	}
	return &resp.Config, nil
}

func (q *queries) GetEnvironments(ctx context.Context, organisationId uuid.UUID, codes []string, page int, pageSize int) (*GetEnvironmentsEnvironmentsPagedEnvironmentsOutput, error) {
	resp, err := GetEnvironments(ctx, q.client, organisationId, codes, page, pageSize)
	if err != nil {
		return nil, err
	}
	return &resp.Environments, nil
}

func (q *queries) GetMemberOrganisations(ctx context.Context, page int, pageSize int) (*GetMemberOrganisationsMemberOrganisationsPagedOrganisationsOutput, error) {
	resp, err := GetMemberOrganisations(ctx, q.client, q.userId, page, pageSize)
	if err != nil {
		return nil, err
	}
	return &resp.MemberOrganisations, nil
}

func (q *queries) NewDeployment(ctx context.Context, organisationId uuid.UUID, environmentId uuid.UUID, configId uuid.UUID, configRevisionId uuid.UUID, revisionId uuid.UUID) (uuid.UUID, error) {
	id := uuid.New()
	resp, err := NewDeployment(ctx, q.client, organisationId, id, environmentId, configId, configRevisionId, revisionId)
	if err != nil {
		return uuid.Nil, err
	}
	return resp.NewDeployment, nil
}

func (q *queries) CreateContainerRepository(ctx context.Context, organisationId uuid.UUID, configId uuid.UUID, code string) (uuid.UUID, error) {
	id := uuid.New()
	resp, err := CreateContainerRepository(ctx, q.client, organisationId, id, configId, code)
	if err != nil {
		return uuid.Nil, err
	}
	return resp.CreateContainerRepository, nil
}

func (q *queries) DeleteContainerRepository(ctx context.Context, organisationId uuid.UUID, id uuid.UUID) (uuid.UUID, error) {
	resp, err := DeleteContainerRepository(ctx, q.client, organisationId, id)
	if err != nil {
		return uuid.Nil, err
	}
	return resp.DeleteContainerRepository, nil
}

func (q *queries) LoginContainerRepository(ctx context.Context, organisationId uuid.UUID, id uuid.UUID) (*AuthContainerRepository, error) {
	resp, err := LoginContainerRepository(ctx, q.client, organisationId, id)
	if err != nil {
		return nil, err
	}
	return &resp.LoginContainerRepository, nil
}

func New[T any](ctx context.Context, cfg *config.NoOps[T]) (Queries, error) {
	httpClient, err := cfg.NewHttpClient(ctx)
	if err != nil {
		return nil, err
	}

	userId, err := cfg.GetUserId()
	if err != nil {
		return nil, err
	}

	orgCode, err := cfg.GetOrganisationCode()
	if err != nil {
		return nil, err
	}

	client := graphql.NewClient(cfg.Api.GraphQL, httpClient)
	return &queries{
		organisation: orgCode,
		userId:       userId,
		client:       client,
	}, nil
}
