package config

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/zitadel/oidc/v2/pkg/oidc"
)

type HomeConfig struct {
	Path string `default:"~/.config/no_ops"`
}

type ApiConfig struct {
	GraphQL string `default:"https://api.getnoops.com/graphql"`
	Token   string `default:""`
}

type AuthConfig struct {
	Issuer   string   `default:"https://account.getnoops.com"`
	ClientId string   `default:"ops"`
	Scopes   []string `default:"openid,profile,email,groups,offline_access"`
}

type LogConfig struct {
	Level string `default:"info"`
}

type Styles struct {
	Title lipgloss.Style
	Desc  lipgloss.Style
	Url   lipgloss.Style
}

type GlobalConfig struct {
	Organisation string `mapstructure:"organisation"`
	Format       string `mapstructure:"format" default:"table"`
}

type Config[C any] struct {
	Organisation string `mapstructure:"organisation"`

	Token   *oidc.Tokens[*oidc.IDTokenClaims]
	Command C
	Global  GlobalConfig
	Home    HomeConfig
	Api     ApiConfig
	Auth    AuthConfig
	Log     LogConfig
}
