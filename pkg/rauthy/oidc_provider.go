package rauthy

import (
	"context"
	"fmt"
)

type OidcProvider struct {
	AdminClaimPath        string `json:"admin_claim_path"`
	AdminClaimValue       string `json:"admin_claim_value"`
	AuthorizationEndpoint string `json:"authorization_endpoint"`
	AutoLink              bool   `json:"auto_link"`
	AutoOnboarding        bool   `json:"auto_onboarding"`
	ClientId              string `json:"client_id"`
	ClientSecret          string `json:"client_secret"`
	ClientSecretBasic     bool   `json:"client_secret_basic"`
	ClientSecretPost      bool   `json:"client_secret_post"`
	Enabled               bool   `json:"enabled"`
	Id                    string `json:"id"`
	Issuer                string `json:"issuer"`
	JwksEndpoint          string `json:"jwks_endpoint"`
	MfaClaimPath          string `json:"mfa_claim_path"`
	MfaClaimValue         string `json:"mfa_claim_value"`
	Name                  string `json:"name"`
	Scope                 string `json:"scope"`
	TokenEndpoint         string `json:"token_endpoint"`
	Typ                   string `json:"typ"`
	UsePkce               bool   `json:"use_pkce"`
	UserinfoEndpoint      string `json:"userinfo_endpoint"`
}

func (c *Client) CreateOidcProvider(ctx context.Context, provider *OidcProvider) (*OidcProvider, error) {
	var newProvider OidcProvider
	_, err := c.Request(ctx, "POST", "/providers/create", &provider, &newProvider)

	if err != nil {
		return nil, err
	}

	return &newProvider, nil
}

func (c *Client) GetOidcProvider(ctx context.Context, id string) (*OidcProvider, error) {
	var providers []OidcProvider
	_, err := c.Request(ctx, "GET", "/providers/", nil, &providers)

	if err != nil {
		return nil, err
	}

	var provider OidcProvider

	for _, p := range providers {
		if p.Id == id {
			provider = p
			break
		}
	}

	if provider == (OidcProvider{}) {
		return nil, fmt.Errorf("no provider found with id %s", id)
	}

	return &provider, nil
}

func (c *Client) UpdateOidcProvider(ctx context.Context, id string, provider *OidcProvider) (*OidcProvider, error) {
	var updatedProvider OidcProvider
	_, err := c.Request(ctx, "PUT", fmt.Sprintf("/providers/%s", id), &provider, &updatedProvider)

	if err != nil {
		return nil, err
	}

	return &updatedProvider, nil
}

func (c *Client) DeleteOidcProvider(ctx context.Context, id string) error {
	_, err := c.Request(ctx, "DELETE", fmt.Sprintf("/providers/%s", id), nil, nil)

	if err != nil {
		return err
	}

	return nil
}
