package models

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type UserBook struct {
	ID         int64     `json:"-"`
	CreatedAt  time.Time `json:"-"`
	ReadAt     time.Time `json:"read_at"`
	ReviewedAt time.Time `json:"reviewed_at"`
	BookID     int64     `json:"book_id"`
	UserID     int64     `json:"user_id"`
	Read       bool      `json:"read"`
	Rating     float32   `json:"rating"`
	ReviewBody string    `json:"review_body"`
	Version    int32     `json:"-"`
}

type UserBookModel struct {
	DB *sql.DB
}

func (ub UserBookModel) Insert(userBook *UserBook) error {
	query := `
    INSERT INTO usersBooksRelation (bookId, userId, read, rating, reviewBody, date_read, reviewed_at)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
    RETURNING id, added_at`

	args := []any{userBook.BookID, userBook.UserID, userBook.Read, userBook.Rating, userBook.ReviewBody, userBook.ReadAt, userBook.ReviewedAt}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return ub.DB.QueryRowContext(ctx, query, args...).Scan(&userBook.ID, &userBook.CreatedAt)
}

func (ub UserBookModel) Update(userBook *UserBook) error {
	query := `
    UPDATE usersBooksRelation 
    SET read = $1, rating = $2, reviewBody = $2, date_read = $3, date_reviewed_at = $4, version = version + 1
    WHERE id = $5 AND version = $6
    RETURNING version`

	args := []any{
		userBook.Read, userBook.Rating, userBook.ReviewBody, userBook.ReadAt, userBook.ReviewedAt, userBook.Version, userBook.ID, userBook.Version,
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
