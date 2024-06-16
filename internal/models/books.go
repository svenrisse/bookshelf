package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/svenrisse/bookshelf/internal/validator"
)

// Book model info
// @Description Book information
type Book struct {
	ID        int64     `json:"-"      example:"5"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"  example:"The Hobbit"`
	Author    string    `json:"author" example:"J.R.R. Tolkien"`
	Year      int32     `json:"year"   example:"1937"`
	Pages     int32     `json:"pages"  example:"320"`
	Genres    []string  `json:"genres" example:"Fantasy,Epic,Children's literature"`
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
	v.Check(len(book.Genres) <= 10, "genres", "must not contain more than 10 genres")
	v.Check(validator.Unique(book.Genres), "genres", "must not contain duplicate values")
}

type BookModel struct {
	DB *sql.DB
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

func (b BookModel) Get(id int64) (*Book, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
    SELECT id, created_at, title, author, year, pages, genres, version 
    FROM books
    WHERE id = $1`

	var book Book

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := b.DB.QueryRowContext(ctx, query, id).
		Scan(&book.ID, &book.CreatedAt, &book.Title, &book.Author, &book.Year, &book.Pages, &book.Genres, &book.Version)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return &book, nil
}

func (b BookModel) Update(book *Book) error {
	query := `
    UPDATE books
    SET title = $1, author = $2, year = $3, pages = $4, genres = $5, version = version + 1
    WHERE id = $6 AND version = $7
    RETURNING version`

	args := []any{
		book.Title,
		book.Author,
		book.Year,
		book.Pages,
		pq.Array(book.Genres),
		book.ID,
		book.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := b.DB.QueryRowContext(ctx, query, args...).Scan(&book.Version)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrEditConflict
		}
		return err
	}

	return nil
}

func (b BookModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `DELETE FROM books
    WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := b.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (b BookModel) ListBooks(
	title string,
	genres []string,
	filters Filters,
) ([]*Book, Metadata, error) {
	query := fmt.Sprintf(`
    SELECT count(*) OVER(), id, created_at, title, author, year, pages, genres, version
    FROM books
    WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
    AND (genres @> $2 OR $2 = '{}')
    ORDER %s %s, id ASC
    LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{title, pq.Array(genres), filters.limit(), filters.offset()}

	rows, err := b.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	books := []*Book{}

	for rows.Next() {
		var book Book

		err := rows.Scan(
			&totalRecords,
			&book.ID,
			&book.CreatedAt,
			&book.Title,
			&book.Author,
			&book.Year,
			&book.Pages,
			pq.Array(&book.Genres),
			&book.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		books = append(books, &book)
	}

	if err = rows.Close(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return books, metadata, nil
}
