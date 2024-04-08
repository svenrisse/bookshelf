package main

import (
	"net/http"
	"testing"

	"github.com/svenrisse/bookshelf/internal/assert"
)

func TestHealthCheck(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, _ := ts.get(t, "/v1/healthcheck")

	assert.Equal(t, code, http.StatusOK)
}
