package auth

import (
	"net/http"
)

// TokenInterceptorTransport implementing RoundTripper interface, using
// this interceptor to add Authorization token header to the request
type TokenInterceptorTransport struct {
	Transport http.RoundTripper
}

func (ti *TokenInterceptorTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	tokens, err := VerifyTokenAndReturn()
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
	return ti.Transport.RoundTrip(req)
}
