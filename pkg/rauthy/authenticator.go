package rauthy

import (
	"fmt"
	"net/http"
)

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
	req.Header.Set("Authorization", fmt.Sprintf("API-Key %s", a.apiKey))
	return nil
}
