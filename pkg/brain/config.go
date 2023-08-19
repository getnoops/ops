package brain

import (
	"github.com/getnoops/ops/cmd/auth"
	"net/http"
	"time"
)

// The global accessible initialised client
var Client *ClientWithResponses

// TokenInterceptorTransport implementing RoundTripper interface, using
// this interceptor to add Authorization token header to the request
type TokenInterceptorTransport struct {
	Transport http.RoundTripper
}

func (ti *TokenInterceptorTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	tokens, err := auth.VerifyTokenAndReturn()
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
	return ti.Transport.RoundTrip(req)
}

func InitClient(serverUrl string) error {
	c := http.Client{
		Transport: &TokenInterceptorTransport{
			Transport: http.DefaultTransport,
		},
		Timeout: time.Duration(90) * time.Second,
	}
	newClient, err := NewClientWithResponses(serverUrl, WithHTTPClient(&c))

	if err != nil {
		return err
	}

	Client = newClient
	return nil
}
