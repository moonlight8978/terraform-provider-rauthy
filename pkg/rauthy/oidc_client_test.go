package rauthy_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/moonlight8978/terraform-provider-rauthy/pkg/rauthy"
	"github.com/stretchr/testify/assert"
)

var oidcClientResponse = `{
			"id": "rauthy",
			"name": "Rauthy",
			"enabled": true,
			"confidential": false,
			"redirect_uris": [
				"https://localhost:8443/auth/v1/oidc/callback"
			],
			"flows_enabled": [
				"authorization_code"
			],
			"access_token_alg": "EdDSA",
			"id_token_alg": "EdDSA",
			"auth_code_lifetime": 10,
			"access_token_lifetime": 10,
			"scopes": [
				"openid"
			],
			"default_scopes": [
				"openid"
			],
			"challenges": [
				"S256"
			],
			"force_mfa": false,
			"client_uri": "https://localhost:8443",
			"contacts": [
				"admin@localhost"
			]
		}`

func TestGetOidcClient(t *testing.T) {
	ts := CreateServer(oidcClientResponse, http.StatusOK)
	defer ts.Close()

	client := rauthy.NewClient(ts.URL, false, rauthy.NewApiKeyAuthenticator("supersecret"))

	oidcClient, err := client.GetOidcClient(context.Background(), "rauthy")

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, oidcClient.Id, "rauthy")
}
