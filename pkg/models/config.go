package models

import (
	"github.com/getnoops/ops/pkg/queries"
)

type Config struct {
	Code         string                      `json:"code"`
	Version      string                      `json:"version"`
	Name         string                      `json:"name"`
	Class        string                      `json:"class"`
	State        string                      `json:"state"`
	Repositories map[string]ConfigRepository `json:"repositories"`
	Resources    map[string]Resource         `json:"resources"`
}

type ConfigRepository struct {
	Code          string `json:"code"`
	State         string `json:"state"`
	RepositoryUri string `json:"repository_uri"`
}

type Resource struct {
	Code string                 `json:"code"`
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

func ToConfig(cfg *queries.Config) *Config {
	repos := map[string]ConfigRepository{}
	for _, repo := range cfg.ContainerRepositories {
		repos[repo.Code] = *ToConfigRepository(repo)
	}

	resources := map[string]Resource{}
	for _, resource := range cfg.Resources {
		resources[resource.Code] = Resource{
			Code: resource.Code,
			Type: string(resource.Type),
			Data: resource.Data,
		}
	}

	return &Config{
		Code:         cfg.Code,
		Version:      cfg.Version_number,
		Name:         cfg.Name,
		Class:        string(cfg.Class),
		State:        string(cfg.State),
		Repositories: repos,
		Resources:    resources,
	}
}

func ToConfigRepository(repo *queries.ContainerRepositoryItem) *ConfigRepository {
	outputs := map[string]string{}
	for _, output := range repo.Stack.Outputs {
		outputs[output.Output_key] = output.Output_value
	}

	return &ConfigRepository{
		Code:          repo.Code,
		State:         string(repo.State),
		RepositoryUri: outputs["RepositoryUri"],
	}
}
