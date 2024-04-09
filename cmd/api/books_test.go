package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/svenrisse/bookshelf/internal/assert"
	"github.com/svenrisse/bookshelf/internal/models"
	"github.com/svenrisse/bookshelf/internal/validator"
)

var validBook = models.Book{
	Title:  "The Hobbit",
	Author: "J.R.R. Tolkien",
	Year:   1932,
	Pages:  300,
	Genres: []string{"Fantasy", "Children's Literature"},
}

func TestValidateBookHandler(t *testing.T) {
	tests := []struct {
		book      models.Book
		name      string
		wantError map[string]string
	}{
		{
			name:      "valid Book",
			book:      validBook,
			wantError: nil,
		},
		{
			name: "invalid Title", book: models.Book{
				Title:  "",
				Author: validBook.Author,
				Year:   validBook.Year,
				Pages:  validBook.Pages,
				Genres: validBook.Genres,
			}, wantError: map[string]string{"title": "must be provided"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := validator.New()
			models.ValidateBook(v, &tt.book)
			assert.DeepEqual(t, tt.wantError, v.Errors)
		})
	}
}

func TestCreateBookHandler(t *testing.T) {
	app := newTestApplication(t)

	w := httptest.NewRecorder()
	tests := []struct {
		book     models.Book
		name     string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid Request",
			book:     validBook,
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
