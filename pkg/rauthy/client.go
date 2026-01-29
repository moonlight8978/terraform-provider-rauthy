package rauthy

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/google/go-querystring/query"
)

type Client struct {
	client        *http.Client
	authenticator Authenticator
	endpoint      string
}

func NewClient(endpoint string, insecure bool, authenticator Authenticator) *Client {
	httpClient := &http.Client{}

	if insecure {
		httpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}

	return &Client{
		httpClient,
		authenticator,
		endpoint,
	}
}

func (c *Client) Request(ctx context.Context, method, path string, payload, responseBody any) (*http.Response, error) {
	var body io.Reader

	switch method {
	case http.MethodPut, http.MethodPost, http.MethodDelete:
		jsonBody, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("Failed to encode JSON body %s %s - Reason: %w", method, path, err)
		}

		body = bytes.NewBuffer(jsonBody)

	default:
		qs, err := query.Values(payload)
		if err != nil {
			return nil, fmt.Errorf("Failed to encode query string %s %s - Reason: %w", method, path, err)
		}

		if encodedQs := qs.Encode(); encodedQs != "" {
			path = fmt.Sprintf("%s?%s", path, encodedQs)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, fmt.Sprintf("%s/%s/%s", c.endpoint, "auth/v1", strings.TrimLeft(path, "/")), body)

	if err != nil {
		return nil, fmt.Errorf("Failed to create request %s %s - Reason: %w", method, path, err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	if err := c.authenticator.Authenticate(req); err != nil {
		return nil, fmt.Errorf("Failed to authenticate request %s %s - Reason: %w", method, path, err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to execute request %s %s - Reason: %w", method, path, err)
	}

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Failed to execute request %s %s - Status Code: %d - Reason: %s", method, path, resp.StatusCode, string(body))
	}

	if responseBody != nil {
		err = json.NewDecoder(resp.Body).Decode(responseBody)
		if err != nil {
			return nil, fmt.Errorf("Failed to decode JSON response %s %s - Reason: %w", method, path, err)
		}
	}

	return resp, nil
}
