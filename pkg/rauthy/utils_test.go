package rauthy_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

func CreateServer(resp string) *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, resp)
	}))

	return ts
}
