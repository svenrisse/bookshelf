package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/svenrisse/bookshelf/internal/assert"
	"github.com/svenrisse/bookshelf/internal/models"
)

func TestCreateBookHandler(t *testing.T) {
	app := newTestApplication(t)

	w := httptest.NewRecorder()

	validBook := models.Book{
		Title:  "The Hobbit",
		Author: "J.R.R. Tolkien",
		Year:   1932,
		Pages:  300,
		Genres: []string{"Fantasy", "Children's Literature"},
	}

	tests := []struct {
		book     models.Book
		name     string
		wantCode int
		wantBody string
	}{
		{
			name: "Valid Request",
			book: models.Book{
				Title:  validBook.Title,
				Author: validBook.Author,
				Year:   validBook.Year,
				Pages:  validBook.Pages,
				Genres: validBook.Genres,
			},
			wantCode: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := json.Marshal(map[string]any{"title": tt.book.Title, "author": tt.book.Author, "year": tt.book.Year, "pages": tt.book.Pages, "genres": tt.book.Genres})
			if err != nil {
				t.Fatal(err)
			}
			req := httptest.NewRequest(http.MethodPost, "/v1/books", bytes.NewBuffer(b))
			app.createBookHandler(w, req)

			assert.Equal(t, w.Result().StatusCode, tt.wantCode)
		})
	}
}
