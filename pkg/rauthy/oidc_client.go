package rauthy

import (
	"context"
	"fmt"
	"net/http"
)

type OidcClient struct {
	ID                  string   `json:"id"`
	Name                string   `json:"name"`
	Enabled             bool     `json:"enabled"`
	Confidential        bool     `json:"confidential"`
	RedirectUris        []string `json:"redirect_uris"`
	FlowsEnabled        []string `json:"flows_enabled"`
	AccessTokenAlg      string   `json:"access_token_alg"`
	IdTokenAlg          string   `json:"id_token_alg"`
	AuthCodeLifetime    int64    `json:"auth_code_lifetime"`
	AccessTokenLifetime int64    `json:"access_token_lifetime"`
	Scopes              []string `json:"scopes"`
	DefaultScopes       []string `json:"default_scopes"`
	Challenges          []string `json:"challenges"`
	ForceMfa            bool     `json:"force_mfa"`
	ClientUri           string   `json:"client_uri"`
	Contacts            []string `json:"contacts"`
	PostLogoutUri       []string `json:"post_logout_redirect_uris,omitempty"`
}

type CreateOidcClientPayload struct {
	ID            string   `json:"id"`
	Confidential  bool     `json:"confidential"`
	Name          string   `json:"name"`
	RedirectUris  []string `json:"redirect_uris"`
	PostLogoutUri []string `json:"post_logout_redirect_uris,omitempty"`
}

func (c *OidcClient) ToCreatePayload() *CreateOidcClientPayload {
	return &CreateOidcClientPayload{
		ID:            c.ID,
		Confidential:  c.Confidential,
		Name:          c.Name,
		RedirectUris:  c.RedirectUris,
		PostLogoutUri: c.PostLogoutUri,
	}
}

func (c *Client) GetOidcClient(ctx context.Context, id string) (*OidcClient, error) {
	var oidcClient OidcClient

	if _, err := c.Request(ctx, http.MethodGet, fmt.Sprintf("/clients/%s", id), nil, &oidcClient); err != nil {
		return nil, err
	}

	return &oidcClient, nil
}

func (c *Client) CreateOidcClient(ctx context.Context, oidcClient *CreateOidcClientPayload) (*OidcClient, error) {
	var createdOidcClient OidcClient

	if _, err := c.Request(ctx, http.MethodPost, "/clients", oidcClient, &createdOidcClient); err != nil {
		return nil, err
	}

	return &createdOidcClient, nil
}

func (c *Client) UpdateOidcClient(ctx context.Context, id string, oidcClient *OidcClient) (*OidcClient, error) {
	var updatedOidcClient OidcClient

	if _, err := c.Request(ctx, http.MethodPut, fmt.Sprintf("/clients/%s", id), oidcClient, &updatedOidcClient); err != nil {
		return nil, err
	}

	return &updatedOidcClient, nil
}

func (c *Client) DeleteOidcClient(ctx context.Context, id string) error {
	if _, err := c.Request(ctx, http.MethodDelete, fmt.Sprintf("/clients/%s", id), nil, nil); err != nil {
		return err
	}

	return nil
}
