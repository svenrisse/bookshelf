package models

import (
	"testing"
	"time"

	"github.com/svenrisse/bookshelf/internal/assert"
)

func TestUserBookModel_Insert(t *testing.T) {
	validUserBook := UserBook{
		BookID:     1,
		UserID:     1,
		Read:       true,
		Rating:     4.5,
		ReviewBody: "Very good book!",
		ReadAt:     time.Date(2024, 04, 12, 14, 30, 00, 0, time.UTC),
		ReviewedAt: time.Date(2024, 04, 14, 18, 00, 00, 0, time.UTC),
	}

	tests := []struct {
		name     string
		userBook UserBook
		wantErr  string
	}{
		{name: "valid", userBook: validUserBook, wantErr: ""},
		{name: "Want to read", userBook: UserBook{
			BookID:     validUserBook.BookID,
			UserID:     validUserBook.BookID,
			Read:       false,
			Rating:     0,
			ReviewBody: "",
		}, wantErr: ""},
		{name: "Non-existent userID", userBook: UserBook{
			BookID:     validUserBook.BookID,
			UserID:     9,
			Read:       validUserBook.Read,
			Rating:     validUserBook.Rating,
			ReviewBody: validUserBook.ReviewBody,
			ReadAt:     validUserBook.ReadAt,
		}, wantErr: "usersbooksrelation_userid_fkey"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := NewTestDB(t)

			ub := UserBookModel{db}

			err := ub.Insert(&tt.userBook)
			if len(tt.wantErr) == 0 {
				assert.NilError(t, err)
				return
			}

			assert.StringContains(t, err.Error(), tt.wantErr)
		})
	}
}
