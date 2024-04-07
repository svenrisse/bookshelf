package mocks

import (
	"time"

	"github.com/svenrisse/bookshelf/internal/models"
)

var mockBook = models.Book{
	ID:        1,
	Title:     "The Hobbit",
	Author:    "J.R.R. Tolkien",
	Year:      1973,
	Pages:     320,
	Genres:    []string{"Fantasy", "Epic", "Children's literature"},
	CreatedAt: time.Now(),
	Version:   1,
}

type BookModel struct{}

func (m *BookModel) Insert(book *models.Book) error {
	return nil
}
