package models

import (
	"context"
	"database/sql"
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
