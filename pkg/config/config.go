package config

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/99designs/keyring"
	"github.com/charmbracelet/lipgloss"
	"github.com/getnoops/ops/pkg/util"
	"github.com/google/uuid"
	"github.com/mcuadros/go-defaults"
	"github.com/zitadel/oidc/v2/pkg/client/rp"
	"github.com/zitadel/oidc/v2/pkg/oidc"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v3"

	"github.com/spf13/viper"
)

var (
	TokenKey         = "token"
	SettingsFilename = "settings.yaml"

	ErrNoOrganisation = errors.New("no organisation set")
)

func openSettings(homePath string) (*os.File, error) {
	settingsPath, err := util.ResolvePath(path.Join(homePath, SettingsFilename))
	if err != nil {
		return nil, err
	}
	baseDir := filepath.Dir(settingsPath)
	if err := os.MkdirAll(baseDir, 0700); err != nil {
		return nil, err
	}
	return os.OpenFile(settingsPath, os.O_RDWR|os.O_CREATE, 0600)
}

type NoOps[C any] struct {
	Config[C]

	writerStderr *os.File
	writerStdout *os.File
	keyring      keyring.Keyring

	Styles Styles
}

func (c *NoOps[C]) StoreToken(token *oidc.Tokens[*oidc.IDTokenClaims]) error {
	wrap := struct {
		Token *oidc.Tokens[*oidc.IDTokenClaims]
	}{
		Token: token,
	}

	raw, err := yaml.Marshal(wrap)
	if err != nil {
		return err
	}

	return c.keyring.Set(keyring.Item{
		Key:  TokenKey,
		Data: raw,
	})
}

func (c *NoOps[C]) StoreSettings(settings map[string]string) error {
	file, err := openSettings(c.Home.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	file.Truncate(0)
	file.Seek(0, 0)

	raw, err := yaml.Marshal(settings)
	if err != nil {
		return err
	}

	if _, err := file.Write(raw); err != nil {
		return err
	}
	return nil
}

func (c *NoOps[C]) GetSettings() (map[string]string, error) {
	file, err := openSettings(c.Home.Path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	raw, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	out := map[string]string{}
	yaml.Unmarshal(raw, &out)
	return out, nil
}

func (c *NoOps[C]) NewRelyingPartyOIDC(ctx context.Context, redirectUri string) (rp.RelyingParty, error) {
	return rp.NewRelyingPartyOIDC(c.Auth.Issuer, c.Auth.ClientId, "", redirectUri, c.Auth.Scopes, rp.WithPKCE(nil))
}

func (c *NoOps[C]) NewHttpClient(ctx context.Context) (*http.Client, error) {
	if c.Token == nil {
		return nil, fmt.Errorf("no token found, please login")
	}

	provider, err := c.NewRelyingPartyOIDC(ctx, "")
	if err != nil {
		c.WriteStderr("failed to create provider")
		return nil, err
	}

	_, verifyErr := rp.VerifyTokens[*oidc.IDTokenClaims](ctx, c.Token.AccessToken, c.Token.IDToken, provider.IDTokenVerifier())
	if errors.Is(verifyErr, oidc.ErrExpired) || errors.Is(verifyErr, oidc.ErrSignatureInvalid) {
		newToken, err := rp.RefreshAccessToken(provider, c.Token.RefreshToken, "", "")
		if err != nil {
			c.WriteStderr("failed to refresh token")
			return nil, err
		}

		c.Token.Token = newToken

		if err := c.StoreToken(c.Token); err != nil {
			c.WriteStderr("failed to store token")
			return nil, err
		}
	}

	return oauth2.NewClient(ctx, oauth2.StaticTokenSource(c.Token.Token)), nil
}

func (c *NoOps[C]) GetUserId() (uuid.UUID, error) {
	if c.Token == nil {
		return uuid.Nil, fmt.Errorf("no token found, please login")
	}

	if c.Token.IDTokenClaims == nil {
		return uuid.Nil, fmt.Errorf("no id token claims found")
	}

	subject := c.Token.IDTokenClaims.Subject
	if subject == "" {
		return uuid.Nil, fmt.Errorf("no subject found in token")
	}

	id, err := uuid.Parse(subject)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (c *NoOps[C]) WriteStderr(out string) {
	c.writerStderr.Write([]byte(out))
}

func (c *NoOps[C]) WriteStdout(out string) {
	c.writerStdout.Write([]byte(out))
}

func New[C any](ctx context.Context, v *viper.Viper) (*NoOps[C], error) {
	var config Config[C]
	defaults.SetDefaults(&config)
	if err := v.Unmarshal(&config); err != nil {
		return nil, err
	}

	// read in the auth.
	ring, err := keyring.Open(keyring.Config{
		ServiceName:      "No_Ops",
		FileDir:          config.Home.Path,
		FilePasswordFunc: keyring.FixedStringPrompt("no_ops"),
	})
	if err != nil {
		return nil, err
	}

	token, err := ring.Get(TokenKey)
	if err != nil && !errors.Is(err, keyring.ErrKeyNotFound) {
		return nil, err
	}
	reader := bytes.NewReader(token.Data)
	if err := v.MergeConfig(reader); err != nil {
		return nil, err
	}

	settings, err := openSettings(config.Home.Path)
	defer settings.Close()

	if err := v.MergeConfig(settings); err != nil {
		return nil, err
	}

	// redo it.
	if err := v.Unmarshal(&config); err != nil {
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
		Config:       config,
		writerStderr: os.Stderr,
		writerStdout: os.Stdout,
		keyring:      ring,
		Styles: Styles{
			Title: titleStyle,
			Desc:  descStyle,
			Url:   urlStyle,
		},
	}, nil
}
