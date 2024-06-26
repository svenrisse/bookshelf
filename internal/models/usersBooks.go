package models

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/svenrisse/bookshelf/internal/validator"
)

type UserBook struct {
	ID         int64     `json:"-"`
	BookID     int64     `json:"book_id"`
	UserID     int64     `json:"user_id"`
	Read       bool      `json:"read"`
	Rating     float32   `json:"rating"`
	ReviewBody string    `json:"review_body"`
	CreatedAt  time.Time `json:"-"`
	ReadAt     time.Time `json:"read_at"`
	ReviewedAt time.Time `json:"reviewed_at"`
	Version    int32     `json:"-"`
}

type UserBookModel struct {
	DB *sql.DB
}

func ValidateUserBook(v *validator.Validator, userBook *UserBook) {
	v.Check(userBook.UserID != 0, "UserID", "must be provided")
	v.Check(userBook.UserID > 0, "UserID", "must be a positive integer")

	v.Check(userBook.BookID != 0, "BookID", "must be provided")
	v.Check(userBook.BookID > 0, "BookID", "must be a positive integer")

	if len(userBook.ReviewBody) != 0 {
		v.Check(len(userBook.ReviewBody) <= 5000, "reviewBody", "must be less than 5000 characters")
		v.Check(userBook.Rating != 0, "rating", "if given reviewBody, rating must be provided")
	}

	if !userBook.ReadAt.IsZero() {
		v.Check(userBook.ReadAt.Year() >= 1900, "ReadAt-Year", "must be greater than 1900")
		v.Check(userBook.ReadAt.Compare(time.Now()) <= 0, "ReadAt", "must not be in the future")
	}

	if !userBook.ReviewedAt.IsZero() {
		v.Check(userBook.ReviewedAt.Year() >= 1900, "ReviewedAt-Year", "must be greater than 1900")
		v.Check(
			userBook.ReviewedAt.Compare(time.Now()) <= 0,
			"ReviewedAt",
			"must not be in the future",
		)
	}
}

func (ub UserBookModel) Insert(userBook *UserBook) error {
	query := `
    INSERT INTO usersBooksRelation (bookId, userId, read, rating, reviewBody, read_at, reviewed_at)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
    RETURNING id, added_at`

	args := []any{
		userBook.BookID,
		userBook.UserID,
		userBook.Read,
		userBook.Rating,
		userBook.ReviewBody,
		userBook.ReadAt,
		userBook.ReviewedAt,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return ub.DB.QueryRowContext(ctx, query, args...).Scan(&userBook.ID, &userBook.CreatedAt)
}

func (ub UserBookModel) Get(id int64) (*UserBook, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `
    SELECT id, bookId, userId, read, rating, reviewBody, added_at, read_at, reviewed_at, version
    FROM usersBooksRelation
    WHERE id = $1`

	var userBook UserBook

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := ub.DB.QueryRowContext(ctx, query, id).Scan(
		&userBook.ID,
		&userBook.BookID,
		&userBook.UserID,
		&userBook.Read,
		&userBook.Rating,
		&userBook.ReviewBody,
		&userBook.ReadAt,
		&userBook.ReviewedAt,
		&userBook.Version)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return &userBook, nil
}

func (ub UserBookModel) Update(userBook *UserBook) error {
	query := `
    UPDATE usersBooksRelation 
    SET read = $1, rating = $2::REAL, reviewBody = $3, read_at = $4, reviewed_at = $5, version = version + 1
    WHERE id = $6 AND version = $7
    RETURNING version`

	args := []any{
		userBook.Read, userBook.Rating, userBook.ReviewBody, userBook.ReadAt, userBook.ReviewedAt, userBook.ID, userBook.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := ub.DB.QueryRowContext(ctx, query, args...).Scan(&userBook.Version)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrEditConflict
		}
		return err
	}

	return nil
}

func (ub UserBookModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `DELETE FROM usersBooksRelation
    WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := ub.DB.ExecContext(ctx, query, id)
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
