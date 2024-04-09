package models

import (
	"testing"

	"github.com/svenrisse/bookshelf/internal/assert"
)

func TestUserModelInsert(t *testing.T) {
	validUser := User{
		Name:  "Insert Jones",
		Email: "insertTest@example.com",
		Password: password{
			hash: NewTestHashPassword(t, "pa55word"),
		},
		Activated: true,
	}
	tests := []struct {
		name  string
		user  User
		error string
	}{
		{name: "Valid User", user: validUser, error: ""},
		{
			name: "Duplicate Email", user: User{
				Name:      validUser.Name,
				Email:     "alice@example.com",
				Password:  validUser.Password,
				Activated: validUser.Activated,
			}, error: "pq: duplicate key value violates unique constraint \"users_email_key\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := NewTestDB(t)

			m := UserModel{db}

			err := m.Insert(&tt.user)

			if len(tt.error) == 0 {
				assert.NilError(t, err)
				return
			}

			assert.StringContains(t, err.Error(), tt.error)
		})
	}
}

func TestUserModelExists(t *testing.T) {
	tests := []struct {
		name   string
		userID int
		want   bool
	}{
		{
			name:   "Valid ID",
			userID: 1,
			want:   true,
		},
		{
			name:   "Zero ID",
			userID: 0,
			want:   false,
		},
		{
			name:   "Non-existent ID",
			userID: 2,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := NewTestDB(t)

			m := UserModel{db}

			exists, err := m.Exists(tt.userID)

			assert.Equal(t, exists, tt.want)
			assert.NilError(t, err)
		})
	}
}
