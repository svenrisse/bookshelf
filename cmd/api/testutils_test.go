package main

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"

	"github.com/svenrisse/bookshelf/internal/mailer"
	"github.com/svenrisse/bookshelf/internal/models/mocks"
)

func newTestApplication(t *testing.T) *application {
	return &application{
		logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		mailer: mailer.New("", 0, "", "", ""),
		models: mocks.NewMockModels(),
	}
}

type testServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewServer(h)

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	ts.Client().Jar = jar

	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, string) {
	rq, err := http.NewRequest("GET", ts.URL+urlPath, nil)
	if err != nil {
		t.Fatal(err)
	}
	rq.Header.Add("Authorization", "Bearer Q5KJHXE3TJ3BUQRFWYYCAFSJDQ")

	rs, err := ts.Client().Do(rq)
	if err != nil {
		t.Fatal(err)
	}

	return rs.StatusCode, rs.Header, string("")
}

func (ts *testServer) post(t *testing.T, urlPath string, body io.Reader) (int, http.Header, string) {
	rq, err := http.NewRequest(http.MethodPost, ts.URL+urlPath, body)
	if err != nil {
		t.Fatal(err)
	}

	rq.Header.Set("Authorization", "Bearer Q5KJHXE3TJ3BUQRFWYYCAFSJDQ")
	rs, err := ts.Client().Do(rq)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	b, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	b = bytes.TrimSpace(b)
	return rs.StatusCode, rs.Header, string(b)
}
