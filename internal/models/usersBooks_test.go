package models

import (
	"testing"
	"time"

	"github.com/svenrisse/bookshelf/internal/assert"
	"github.com/svenrisse/bookshelf/internal/validator"
)

var validUserBook = UserBook{
	BookID:     1,
	UserID:     1,
	Read:       true,
	Rating:     4.5,
	ReviewBody: "Very good book!",
	ReadAt:     time.Date(2024, 04, 12, 14, 30, 00, 0, time.UTC),
	ReviewedAt: time.Date(2024, 04, 14, 18, 00, 00, 0, time.UTC),
}

func TestUserBookModel_Insert(t *testing.T) {
	tests := []struct {
		name     string
		userBook UserBook
		wantErr  string
	}{
		{
			name:     "valid",
			userBook: validUserBook,
			wantErr:  "",
		},
		{
			name: "Want to read",
			userBook: UserBook{
				BookID: validUserBook.BookID, UserID: validUserBook.BookID, Read: false, Rating: 0, ReviewBody: "",
			},
			wantErr: "",
		},
		{
			name: "Non-existent userID",
			userBook: UserBook{
				BookID:     validUserBook.BookID,
				UserID:     9,
				Read:       validUserBook.Read,
				Rating:     validUserBook.Rating,
				ReviewBody: validUserBook.ReviewBody,
				ReadAt:     validUserBook.ReadAt,
			},
			wantErr: "usersbooksrelation_userid_fkey",
		},
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

func TestUserBookModel_Update(t *testing.T) {
	tests := []struct {
		name    string
		wantErr string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func TestUserBookModel_Delete(t *testing.T) {
	tests := []struct {
		name    string
		id      int64
		wantErr error
	}{
		{name: "Valid Delete", id: 14, wantErr: nil},
		{name: "Invalid Delete", id: 10, wantErr: ErrRecordNotFound},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := NewTestDB(t)

			ub := UserBookModel{db}

			err := ub.Delete(tt.id)

			assert.Equal(t, err, tt.wantErr)
		})
	}
}

func TestValidateUserBook(t *testing.T) {
	tests := []struct {
		name      string
		userBook  UserBook
		wantError map[string]string
	}{
		{
			name:      "Want to read",
			userBook:  UserBook{UserID: validUserBook.UserID, BookID: validUserBook.BookID},
			wantError: nil,
		},
		{
			name:      "Valid UserBook",
			userBook:  validUserBook,
			wantError: nil,
		},
		{
			name: "Missing Rating when reviewing",
			userBook: UserBook{
				UserID:     validUserBook.UserID,
				BookID:     validUserBook.BookID,
				Read:       validUserBook.Read,
				ReviewBody: "This is the review body.",
			},
			wantError: map[string]string{"rating": "if given reviewBody, rating must be provided"},
		},
		{
			name: "Read At in future",
			userBook: UserBook{
				UserID:     validUserBook.UserID,
				BookID:     validUserBook.BookID,
				Read:       validUserBook.Read,
				ReviewBody: validUserBook.ReviewBody,
				Rating:     validUserBook.Rating,
				ReadAt:     time.Now().Add(5 * time.Hour),
			},
			wantError: map[string]string{"ReadAt": "must not be in the future"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := validator.New()

			ValidateUserBook(v, &tt.userBook)
			assert.DeepEqual(t, tt.wantError, v.Errors)
		})
	}
}
