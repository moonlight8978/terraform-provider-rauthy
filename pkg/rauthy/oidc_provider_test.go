package rauthy_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/moonlight8978/terraform-provider-rauthy/pkg/rauthy"
	"github.com/stretchr/testify/assert"
)

var oidcProviderResponse = `{
	"admin_claim_path": "admin",
	"admin_claim_value": "true",
	"authorization_endpoint": "https://accounts.google.com/o/oauth2/v2/auth",
	"auto_link": false,
	"auto_onboarding": false,
	"client_id": "google-client-id",
	"client_secret": "google-client-secret",
	"client_secret_basic": true,
	"client_secret_post": false,
	"enabled": true,
	"id": "google",
	"issuer": "https://accounts.google.com",
	"jwks_endpoint": "https://www.googleapis.com/oauth2/v3/certs",
	"mfa_claim_path": "mfa",
	"mfa_claim_value": "true",
	"name": "Google",
	"scope": "openid profile email",
	"token_endpoint": "https://oauth2.googleapis.com/token",
	"typ": "oidc",
	"use_pkce": true,
	"userinfo_endpoint": "https://openidconnect.googleapis.com/v1/userinfo"
}`

func TestCreateOidcProvider(t *testing.T) {
	ts := CreateServer(oidcProviderResponse, http.StatusOK)
	defer ts.Close()

	client := rauthy.NewClient(ts.URL, false, rauthy.NewApiKeyAuthenticator("supersecret"))

	p := &rauthy.OidcProvider{
		Id:           "google",
		Name:         "Google",
		ClientId:     "google-client-id",
		ClientSecret: "google-client-secret",
	}

	provider, err := client.CreateOidcProvider(context.Background(), p)
	assert.Nil(t, err)
	assert.Equal(t, "google", provider.Id)
	assert.Equal(t, "Google", provider.Name)
}

func TestGetOidcProvider(t *testing.T) {
	ts := CreateServer("["+oidcProviderResponse+"]", http.StatusOK)
	defer ts.Close()

	client := rauthy.NewClient(ts.URL, false, rauthy.NewApiKeyAuthenticator("supersecret"))

	provider, err := client.GetOidcProvider(context.Background(), "google")
	assert.Nil(t, err)
	assert.Equal(t, "google", provider.Id)
	assert.Equal(t, "Google", provider.Name)
}

func TestGetOidcProvider_NotFound(t *testing.T) {
	ts := CreateServer("["+oidcProviderResponse+"]", http.StatusOK)
	defer ts.Close()

	client := rauthy.NewClient(ts.URL, false, rauthy.NewApiKeyAuthenticator("supersecret"))

	_, err := client.GetOidcProvider(context.Background(), "not-exists")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "no provider found with id")
}

func TestUpdateOidcProvider(t *testing.T) {
	ts := CreateServer(oidcProviderResponse, http.StatusOK)
	defer ts.Close()

	client := rauthy.NewClient(ts.URL, false, rauthy.NewApiKeyAuthenticator("supersecret"))

	p := &rauthy.OidcProvider{
		Id:   "google",
		Name: "Google Updated",
	}

	provider, err := client.UpdateOidcProvider(context.Background(), "google", p)
	assert.Nil(t, err)
	// The mock response has Name="Google", so this confirms the response is parsed correctly.
	assert.Equal(t, "Google", provider.Name)
}

func TestDeleteOidcProvider(t *testing.T) {
	ts := CreateServer("", http.StatusOK)
	defer ts.Close()

	client := rauthy.NewClient(ts.URL, false, rauthy.NewApiKeyAuthenticator("supersecret"))

	err := client.DeleteOidcProvider(context.Background(), "google")
	assert.Nil(t, err)
}
