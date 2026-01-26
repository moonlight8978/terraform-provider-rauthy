package rauthy

import "net/http"

type Authenticator interface {
	Authenticate(req *http.Request) error
}

type ApiKeyAuthenticator struct {
	apiKey string
}

func NewApiKeyAuthenticator(apiKey string) *ApiKeyAuthenticator {
	return &ApiKeyAuthenticator{
		apiKey: apiKey,
	}
}

func (a *ApiKeyAuthenticator) Authenticate(req *http.Request) error {
	req.Header.Set("X-API-Key", a.apiKey)
	return nil
}
