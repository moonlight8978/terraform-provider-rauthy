package rauthy

import (
	"crypto/tls"
	"net/http"
)

type Client struct {
	client        *http.Client
	authenticator *Authenticator
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
		&authenticator,
	}
}
