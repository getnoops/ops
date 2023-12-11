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
	GetMemberOrganisations(ctx context.Context, page int, pageSize int) (*GetMemberOrganisationsMemberOrganisationsPagedOrganisationsOutput, error)
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

func (q *queries) GetMemberOrganisations(ctx context.Context, page int, pageSize int) (*GetMemberOrganisationsMemberOrganisationsPagedOrganisationsOutput, error) {
	resp, err := GetMemberOrganisations(ctx, q.client, q.userId, page, pageSize)
	if err != nil {
		return nil, err
	}
	return &resp.MemberOrganisations, nil
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

	client := graphql.NewClient(cfg.Api.GraphQL, httpClient)
	return &queries{
		organisation: cfg.Organisation,
		userId:       userId,
		client:       client,
	}, nil
}
