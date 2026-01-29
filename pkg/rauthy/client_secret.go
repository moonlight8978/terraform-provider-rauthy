package rauthy

import (
	"context"
	"fmt"
	"net/http"
)

type ClientSecret struct {
	Id           string `json:"id"`
	Confidential bool   `json:"confidential"`
	Secret       string `json:"secret"`
}

type ClientSecretRequest struct {
	CacheCurrentHours int `json:"cache_current_hours,omitempty"`
}

func (c *Client) CreateClientSecret(ctx context.Context, clientId string, req *ClientSecretRequest) (*ClientSecret, error) {
	var secret ClientSecret
	var err error

	if req.CacheCurrentHours == 0 {
		_, err = c.Request(ctx, http.MethodPut, fmt.Sprintf("clients/%s/secret", clientId), nil, &secret)
	} else {
		_, err = c.Request(ctx, http.MethodPost, fmt.Sprintf("clients/%s/secret", clientId), &req, &secret)
	}

	if err != nil {
		return nil, err
	}

	return &secret, nil
}

func (c *Client) GetClientSecret(ctx context.Context, clientId string) (*ClientSecret, error) {
	var secret ClientSecret
	_, err := c.Request(ctx, http.MethodPost, fmt.Sprintf("clients/%s/secret", clientId), nil, &secret)

	if err != nil {
		return nil, err
	}

	return &secret, nil
}
