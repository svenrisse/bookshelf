package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/svenrisse/bookshelf/internal/assert"
)

func TestCreateBookHandler(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())

	const (
		validTitle  = "The Hobbit"
		validAuthor = "J.R.R Tolkien"
		validYear   = 1932
		validPages  = 300
	)

	tests := []struct {
		title    string
		author   string
		year     int
		pages    int
		genres   []string
		name     string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid Request",
			wantCode: http.StatusCreated,
			title:    validTitle,
			author:   validAuthor,
			year:     validYear,
			pages:    validPages,
			genres:   []string{"Fantasy", "Children's Literature"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := json.Marshal(map[string]any{"title": tt.title, "author": tt.author, "year": tt.year, "pages": tt.pages, "genres": tt.genres})
			if err != nil {
				t.Fatal(err)
			}
			code, header, body := ts.post(t, "/v1/books", bytes.NewBuffer(b))

			t.Log(code)
			t.Log(header)
			t.Log(body)
			assert.Equal(t, code, tt.wantCode)
		})
	}
}
