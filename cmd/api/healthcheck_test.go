package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/svenrisse/bookshelf/internal/assert"
)

func TestHealthCheck(t *testing.T) {
	app := newTestApplication(t)
	req := httptest.NewRequest(http.MethodGet, "/v1/healthcheck", nil)
	w := httptest.NewRecorder()
	req.Header.Set("Authorization", "Bearer Q5KJHXE3TJ3BUQRFWYYCAFSJDQ")
	app.healthcheckHandler(w, req)

	resp := w.Result()
	assert.Equal(t, resp.StatusCode, http.StatusOK)
}
