package login

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/charmbracelet/lipgloss"
	"github.com/getnoops/ops/pkg/config"

	"github.com/google/uuid"
	"github.com/zitadel/oidc/v2/pkg/client/rp"
	"github.com/zitadel/oidc/v2/pkg/oidc"
)

func codeExchangeHandler[C oidc.IDClaims](callback rp.CodeExchangeCallback[C], provider rp.RelyingParty, state string, codeVerifier string, urlParam ...rp.URLParamOpt) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		if params.Get("error") != "" {
			provider.ErrorHandler()(w, r, params.Get("error"), params.Get("error_description"), state)
			return
		}

		codeOpts := make([]rp.CodeExchangeOpt, len(urlParam))
		for i, p := range urlParam {
			codeOpts[i] = rp.CodeExchangeOpt(p)
		}

		if provider.IsPKCE() {
			codeOpts = append(codeOpts, rp.WithCodeVerifier(codeVerifier))
		}

		tokens, err := rp.CodeExchange[C](r.Context(), params.Get("code"), provider, codeOpts...)
		if err != nil {
			http.Error(w, "failed to exchange token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		callback(w, r, tokens, state, provider)
	}
}

func authRedirect(provider rp.RelyingParty, state string, codeVerifier string, urlParam ...rp.URLParamOpt) string {
	var opts []rp.AuthURLOpt
	for _, p := range urlParam {
		opts = append(opts, rp.AuthURLOpt(p))
	}

	if provider.IsPKCE() {
		codeChallenge := oidc.NewSHACodeChallenge(codeVerifier)
		opts = append(opts, rp.WithCodeChallenge(codeChallenge))
	}

	return rp.AuthURL(state, provider, opts...)
}

type Server interface {
	Shutdown(ctx context.Context) error
}

func NewServer(ctx context.Context, cfg *config.NoOps[Config, string], tokenChan chan *oidc.Tokens[*oidc.IDTokenClaims]) (Server, error) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		cfg.WriteStderr("failed to listen on port")
		return nil, err
	}

	port := l.Addr().(*net.TCPAddr).Port
	redirectUri := fmt.Sprintf("http://localhost:%d/auth/callback", port)

	provider, err := cfg.NewRelyingPartyOIDC(ctx, redirectUri)
	if err != nil {
		cfg.WriteStderr("failed to create provider")
		return nil, err
	}

	callback := func(w http.ResponseWriter, r *http.Request, tokens *oidc.Tokens[*oidc.IDTokenClaims], state string, rp rp.RelyingParty) {
		tokenChan <- tokens

		// todo add a redirect here instead of a message
		msg := "<p><strong>Success!</strong></p>"
		msg = msg + "<p>You are authenticated and can now return to the CLI.</p>"
		w.Write([]byte(msg))
	}

	state := uuid.NewString()
	codeVerifier := base64.RawURLEncoding.EncodeToString([]byte(uuid.New().String()))

	mux := http.NewServeMux()
	mux.Handle("/auth/callback", codeExchangeHandler(callback, provider, state, codeVerifier))

	srv := &http.Server{
		Handler: mux,
	}

	go func() {
		if err := srv.Serve(l); err != nil && !errors.Is(err, http.ErrServerClosed) {
			cfg.WriteStderr("failed to serve")
		}
	}()

	url := authRedirect(provider, state, codeVerifier)

	out := lipgloss.JoinVertical(
		lipgloss.Left,
		cfg.Styles.Title.Render("To authenticate please follow the link below:"),
		cfg.Styles.Url.Render(url),
	)

	cfg.WriteStdout(out)
	return srv, nil
}
