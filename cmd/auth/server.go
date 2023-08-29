package auth

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/exec"
	"runtime"

	"github.com/getnoops/ops/pkg/tokenstore"
	"github.com/spf13/viper"

	"github.com/getnoops/ops/pkg/logging"
	"github.com/google/uuid"
	"github.com/zitadel/oidc/v2/pkg/client/rp"
	"github.com/zitadel/oidc/v2/pkg/oidc"
)

func OpenBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}

func CodeExchangeHandler[C oidc.IDClaims](callback rp.CodeExchangeCallback[C], provider rp.RelyingParty, state string, codeVerifier string, urlParam ...rp.URLParamOpt) http.HandlerFunc {
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
		// if rp.Signer() != nil {
		// 	assertion, err := client.SignedJWTProfileAssertion(rp.OAuthConfig().ClientID, []string{rp.Issuer()}, time.Hour, rp.Signer())
		// 	if err != nil {
		// 		http.Error(w, "failed to build assertion: "+err.Error(), http.StatusUnauthorized)
		// 		return
		// 	}
		// 	codeOpts = append(codeOpts, WithClientAssertionJWT(assertion))
		// }

		tokens, err := rp.CodeExchange[C](r.Context(), params.Get("code"), provider, codeOpts...)
		if err != nil {
			http.Error(w, "failed to exchange token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		callback(w, r, tokens, state, provider)
	}
}

// VerifyTokenAndReturn verifies the token and handles the case of expiration of token
// - it will retrieve the tokens from [User Home Directory]/.no_ops/no_opsconfig file
// - checks whether the token is expired or not, if expired it refreshes and stores the new token
// - returns the token if they are valid
// - returns an error for issues in
//   - retrieving token from file
//   - verifying token
//   - refreshing the access token
//   - storing updated tokens
func VerifyTokenAndReturn() (*tokenstore.Tokens, error) {
	token, err := tokenstore.Retrieve()
	if err != nil {
		return nil, err
	}
	config := MustNewConfig(viper.GetViper())
	options := []rp.Option{
		rp.WithPKCE(nil),
	}

	provider, err := rp.NewRelyingPartyOIDC(config.Auth.Issuer, config.Auth.ClientId, "", "", config.Auth.Scopes, options...)
	logging.OnError(err).Fatal("error creating provider")

	// Don't need claims here for now, may need later
	_, err = rp.VerifyTokens[*oidc.IDTokenClaims](context.Background(), token.AccessToken, token.IDToken, provider.IDTokenVerifier())
	if err != nil {
		if errors.Is(err, oidc.ErrExpired) || errors.Is(err, oidc.ErrSignatureInvalid) {
			newAccessToken, refreshAccessTokenErr := rp.RefreshAccessToken(provider, token.RefreshToken, "", "")
			if refreshAccessTokenErr != nil {
				return nil, refreshAccessTokenErr
			}
			updateTokensErr := tokenstore.UpdateTokens(newAccessToken.AccessToken, newAccessToken.RefreshToken)
			if updateTokensErr != nil {
				return nil, updateTokensErr
			}
			return &tokenstore.Tokens{
				AccessToken:  newAccessToken.AccessToken,
				RefreshToken: newAccessToken.RefreshToken,
				IDToken:      token.IDToken,
				TokenType:    token.TokenType,
			}, nil
		}
		return nil, err
	}
	return token, nil
}

func AuthRedirect(provider rp.RelyingParty, state string, codeVerifier string, urlParam ...rp.URLParamOpt) {
	opts := make([]rp.AuthURLOpt, len(urlParam))

	if provider.IsPKCE() {
		codeChallenge := oidc.NewSHACodeChallenge(codeVerifier)
		opts = append(opts, rp.WithCodeChallenge(codeChallenge))
	}

	url := rp.AuthURL(state, provider, opts...)
	OpenBrowser(url)
	fmt.Printf("URL: %s\n", url)
}

type Server interface {
	Shutdown(ctx context.Context) error
}

func NewServer(ctx context.Context, config *Config, tokenChan chan *oidc.Tokens[*oidc.IDTokenClaims]) (Server, error) {
	options := []rp.Option{
		rp.WithPKCE(nil),
	}

	l, err := net.Listen("tcp", ":0")
	logging.OnError(err).Fatal("error listening on port")
	port := l.Addr().(*net.TCPAddr).Port
	redirectUri := fmt.Sprintf("http://localhost:%d/auth/callback", port)
	state := uuid.NewString()
	codeVerifier := base64.RawURLEncoding.EncodeToString([]byte(uuid.New().String()))

	provider, err := rp.NewRelyingPartyOIDC(config.Auth.Issuer, config.Auth.ClientId, "", redirectUri, config.Auth.Scopes, options...)
	logging.OnError(err).Fatal("error creating provider")

	callback := func(w http.ResponseWriter, r *http.Request, tokens *oidc.Tokens[*oidc.IDTokenClaims], state string, rp rp.RelyingParty) {
		tokenChan <- tokens
		msg := "<p><strong>Success!</strong></p>"
		msg = msg + "<p>You are authenticated and can now return to the CLI.</p>"
		w.Write([]byte(msg))
	}

	mux := http.NewServeMux()
	mux.Handle("/auth/callback", CodeExchangeHandler(callback, provider, state, codeVerifier))

	srv := &http.Server{
		Handler: mux,
	}

	go func() {
		if err := srv.Serve(l); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logging.OnError(err).Fatal("error starting server")
		}
	}()

	AuthRedirect(provider, state, codeVerifier)
	return srv, nil
}
