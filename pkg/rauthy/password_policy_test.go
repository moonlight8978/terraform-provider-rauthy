package rauthy_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/moonlight8978/terraform-provider-rauthy/pkg/rauthy"
	"github.com/stretchr/testify/assert"
)

func TestGetPasswordPolicy(t *testing.T) {
	ts := createServer()
	defer ts.Close()

	client := rauthy.NewClient(ts.URL, false, rauthy.NewApiKeyAuthenticator("supersecret"))

	passwordPolicy, err := client.GetPasswordPolicy(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 8, passwordPolicy.LengthMin)
	assert.Equal(t, 128, passwordPolicy.LengthMax)
	assert.Equal(t, 1, passwordPolicy.IncludeLowerCase)
	assert.Equal(t, 1, passwordPolicy.IncludeUpperCase)
	assert.Equal(t, 1, passwordPolicy.IncludeDigits)
	assert.Equal(t, 180, passwordPolicy.ValidDays)
	assert.Equal(t, 3, passwordPolicy.NotRecentlyUsed)
}

func TestUpdatePasswordPolicy(t *testing.T) {
	ts := createServer()
	defer ts.Close()

	client := rauthy.NewClient(ts.URL, false, rauthy.NewApiKeyAuthenticator("supersecret"))

	passwordPolicy, err := client.UpdatePasswordPolicy(context.Background(), &rauthy.PasswordPolicy{
		LengthMin:        6,
		LengthMax:        128,
		IncludeLowerCase: 1,
		IncludeUpperCase: 2,
		IncludeDigits:    1,
		ValidDays:        180,
		NotRecentlyUsed:  3,
	})
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 8, passwordPolicy.LengthMin)
	assert.Equal(t, 128, passwordPolicy.LengthMax)
	assert.Equal(t, 1, passwordPolicy.IncludeLowerCase)
	assert.Equal(t, 1, passwordPolicy.IncludeUpperCase)
	assert.Equal(t, 1, passwordPolicy.IncludeDigits)
	assert.Equal(t, 180, passwordPolicy.ValidDays)
	assert.Equal(t, 3, passwordPolicy.NotRecentlyUsed)
}

func createServer() *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{
			"length_min": 8,
			"length_max": 128,
			"include_lower_case": 1,
			"include_upper_case": 1,
			"include_digits": 1,
			"valid_days": 180,
			"not_recently_used": 3
		}`)
	}))

	return ts
}
