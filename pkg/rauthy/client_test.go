package rauthy_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/moonlight8978/terraform-provider-rauthy/pkg/rauthy"
	"github.com/stretchr/testify/assert"
)

func TestRequest_Ok(t *testing.T) {
	ts := CreateServer(`{"id": "rauthy"}`, http.StatusOK)
	defer ts.Close()

	client := rauthy.NewClient(ts.URL, false, rauthy.NewApiKeyAuthenticator("supersecret"))

	_, err := client.Request(context.Background(), "GET", "/test", nil, nil)

	assert.NoError(t, err)
}

func TestRequest_BadRequest(t *testing.T) {
	ts := CreateServer(``, http.StatusBadRequest)
	defer ts.Close()

	client := rauthy.NewClient(ts.URL, false, rauthy.NewApiKeyAuthenticator("supersecret"))

	_, err := client.Request(context.Background(), "GET", "/test", nil, nil)

	assert.Error(t, err)
}
