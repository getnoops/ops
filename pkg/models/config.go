package models

import (
	"github.com/getnoops/ops/pkg/queries"
	"github.com/google/uuid"
)

type Config struct {
	Code         string                      `json:"code"`
	Version      string                      `json:"version"`
	Name         string                      `json:"name"`
	Class        string                      `json:"class"`
	State        string                      `json:"state"`
	Registry     *Registry                   `json:"registry"`
	Repositories map[string]ConfigRepository `json:"repositories"`
	Secrets      map[string]Secret           `json:"secrets"`
	Resources    map[string]Resource         `json:"resources"`
}

type Registry struct {
	Username    string `json:"username"`
	RegistryUrl string `json:"registry_url"`
}

type ConfigRepository struct {
	Code          string `json:"code"`
	State         string `json:"state"`
	RepositoryUri string `json:"repository_uri"`
}

type Secret struct {
	Id          uuid.UUID `json:"id"`
	Code        string    `json:"code"`
	State       string    `json:"state"`
	Environment string    `json:"environment"`
}

type Resource struct {
	Code string                 `json:"code"`
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

func ToConfig(cfg *queries.Config) *Config {
	secrets := map[string]Secret{}
	for _, secret := range cfg.Secrets {
		secrets[secret.Code] = *ToSecret(secret)
	}

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

	var registry *Registry
	if cfg.Registry != nil {
		registry = &Registry{
			Username:    cfg.Registry.Username,
			RegistryUrl: cfg.Registry.Registry_url,
		}
	}

	return &Config{
		Code:         cfg.Code,
		Version:      cfg.Version_number,
		Name:         cfg.Name,
		Class:        string(cfg.Class),
		State:        string(cfg.State),
		Repositories: repos,
		Secrets:      secrets,
		Registry:     registry,
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

func ToSecret(secret *queries.SecretItem) *Secret {
	return &Secret{
		Id:          secret.Id,
		Code:        secret.Code,
		State:       string(secret.State),
		Environment: secret.Environment.Code,
	}
}

func T(repo *queries.ContainerRepositoryItem) *ConfigRepository {
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
