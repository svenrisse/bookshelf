package models

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
	"github.com/svenrisse/bookshelf/internal/validator"
)

// Book model info
// @Description Book information
type Book struct {
	ID        int64     `json:"-"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"             example:"The Hobbit"`
	Author    string    `json:"author"            example:"J.R.R. Tolkien"`
	Year      int32     `json:"year,omitempty"    example:"1937"`
	Pages     int32     `json:"pages,omitempty"   example:"320"`
	Genres    []string  `json:"genres,omitempty"  example:"Fantasy,Epic,Children's literature"`
	Version   int32     `json:"-"`
}

func ValidateBook(v *validator.Validator, book *Book) {
	v.Check(book.Title != "", "title", "must be provided")
	v.Check(len(book.Title) <= 500, "title", "must not be more than 500 bytes long")

	v.Check(book.Author != "", "author", "must be provided")
	v.Check(len(book.Author) <= 300, "author", "must not be more than 300 bytes long")

	v.Check(book.Year != 0, "year", "must be provided")
	v.Check(book.Year <= int32(time.Now().Year()), "year", "must not be in the future")

	v.Check(book.Pages != 0, "pages", "must be provided")
	v.Check(book.Pages >= 0, "pages", "must be a positive integer")

	v.Check(book.Genres != nil, "genres", "must be provided")
	v.Check(len(book.Genres) >= 0, "genres", "must contain atleast 1 genre")
	v.Check(len(book.Genres) <= 10, "genres", "must not conatin more than 10 genres")
	v.Check(validator.Unique(book.Genres), "genres", "must not contain duplicate values")
}

type BookModel struct {
	DB *sql.DB
}

type BookModelInterface interface {
	Insert(book *Book) error
}

func (b BookModel) Insert(book *Book) error {
	query := `
    INSERT INTO books (title, author, year, pages, genres)
    VALUES ($1, $2, $3, $4, $5,)
    RETURNING id, created_at, version`

	args := []any{book.Title, book.Author, book.Year, book.Pages, pq.Array(book.Genres)}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return b.DB.QueryRowContext(ctx, query, args...).Scan(&book.ID, &book.CreatedAt, &book.Version)
}
