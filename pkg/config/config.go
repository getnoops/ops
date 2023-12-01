package config

import (
	"context"
	"encoding/json"
	"os"

	"github.com/99designs/keyring"
	"github.com/charmbracelet/lipgloss"
	"github.com/mcuadros/go-defaults"
	"github.com/zitadel/oidc/v2/pkg/oidc"
	"go.uber.org/zap"

	"github.com/spf13/viper"
)

var TokenKey = "token"

type AuthConfig struct {
	Issuer   string   `default:"https://account.getnoops.com"`
	ClientId string   `default:"ops"`
	Scopes   []string `default:"openid,profile,email,groups,offline_access"`
}

type HomeConfig struct {
	Path string `default:"~/.config/no_ops"`
}

type LogConfig struct {
	Level string `default:"info"`
}

type Styles struct {
	Title lipgloss.Style
	Desc  lipgloss.Style
	Url   lipgloss.Style
}

type Config[C any] struct {
	Command C
	Home    HomeConfig
	Auth    AuthConfig
	Log     LogConfig
}

type NoOps[C any] struct {
	Config[C]

	writer  *os.File
	keyring keyring.Keyring

	Styles Styles
	Log    *zap.Logger
}

func (c *NoOps[C]) StoreToken(token *oidc.Tokens[*oidc.IDTokenClaims]) error {
	raw, err := json.Marshal(token)
	if err != nil {
		return err
	}

	return c.keyring.Set(keyring.Item{
		Key:  TokenKey,
		Data: raw,
	})
}

func (c *NoOps[C]) GetToken() (*oidc.Tokens[*oidc.IDTokenClaims], error) {
	value, err := c.keyring.Get(TokenKey)
	if err != nil {
		return nil, err
	}

	var out oidc.Tokens[*oidc.IDTokenClaims]
	if err := json.Unmarshal(value.Data, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *NoOps[C]) Write(out string) {
	c.writer.Write([]byte(out))
}

func New[C any](ctx context.Context, v *viper.Viper) (*NoOps[C], error) {
	var config Config[C]
	defaults.SetDefaults(&config)
	if err := v.Unmarshal(&config); err != nil {
		return nil, err
	}

	lvl, err := zap.ParseAtomicLevel(config.Log.Level)
	if err != nil {
		return nil, err
	}

	logCfg := zap.NewDevelopmentConfig()
	logCfg.Level = lvl

	logger, err := logCfg.Build()
	if err != nil {
		return nil, err
	}

	ring, err := keyring.Open(keyring.Config{
		ServiceName:      "No_Ops",
		FileDir:          config.Home.Path,
		FilePasswordFunc: keyring.FixedStringPrompt("no_ops"),
	})
	if err != nil {
		return nil, err
	}

	special := lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}

	re := lipgloss.NewRenderer(os.Stdout)
	descStyle := re.NewStyle().MarginTop(1)
	urlStyle := re.NewStyle().Foreground(special)
	titleStyle := re.NewStyle().
		MarginLeft(1).
		MarginRight(5).
		Padding(0, 1).
		Italic(true).
		Foreground(lipgloss.Color("#FFF7DB"))

	return &NoOps[C]{
		Config:  config,
		writer:  os.Stdout,
		keyring: ring,
		Log:     logger,
		Styles: Styles{
			Title: titleStyle,
			Desc:  descStyle,
			Url:   urlStyle,
		},
	}, nil
}
