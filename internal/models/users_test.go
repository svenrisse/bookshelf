package models

import (
	"testing"

	"github.com/svenrisse/bookshelf/internal/assert"
)

// func TestValidateUser(t *testing.T) {
// 	testPlaintextPassword := "pa55word"
// 	shortPlaintextPassword := strings.Repeat("1", 6)
// 	longPlaintextPassword := strings.Repeat("1", 75)
// 	validUser := User{
// 		Name:  "Validate Jones",
// 		Email: "validatetest@example.com",
// 		Password: password{
// 			plaintext: &testPlaintextPassword,
// 			hash:      NewTestHashPassword(t, testPlaintextPassword),
// 		},
// 	}
//
// 	tests := []struct {
// 		name      string
// 		user      User
// 		wantError map[string]string
// 	}{
// 		{name: "Valid User", user: validUser, wantError: nil},
// 		{name: "To short Password", user: User{
// 			Name: validUser.Name, Email: validUser.Email, Password: password{plaintext: &shortPlaintextPassword, hash: NewTestHashPassword(t, shortPlaintextPassword)},
// 		}, wantError: map[string]string{"password": "must be at least 8 bytes long"}},
// 		{name: "To long Password", user: User{
// 			Name: validUser.Name, Email: validUser.Email, Password: password{plaintext: &longPlaintextPassword, hash: NewTestHashPassword(t, shortPlaintextPassword)},
// 		}, wantError: map[string]string{"password": "must not be more than 72 bytes long"}},
// 		{name: "Missing Name", user: User{
// 			Name: "", Email: validUser.Email, Password: validUser.Password,
// 		}, wantError: map[string]string{"name": "must be provided"}},
// 		{name: "To long Name", user: User{
// 			Name: strings.Repeat("1", 600), Email: validUser.Email, Password: validUser.Password,
// 		}, wantError: map[string]string{"name": "must not be more than 500 bytes long"}},
// 		{name: "Missing Email", user: User{
// 			Name: validUser.Name, Email: "", Password: validUser.Password,
// 		}, wantError: map[string]string{"email": "must be provided"}},
// 		{name: "Invalid Email", user: User{
// 			Name: validUser.Name, Email: "aliceJones.com", Password: validUser.Password,
// 		}, wantError: map[string]string{"email": "must be a valid email address"}},
// 	}
//
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			v := validator.New()
//
// 			ValidateUser(v, &tt.user)
// 			assert.DeepEqual(t, tt.wantError, v.Errors)
// 		})
// 	}
// }

func TestUserModelInsert(t *testing.T) {
	validUser := User{
		ID:       123456,
		Name:     "Insert Jones",
		Avatar:   "avat",
		Provider: "discord",
	}
	tests := []struct {
		name  string
		user  User
		error string
	}{
		{name: "Valid User Insert", user: validUser, error: ""},
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
		name string
		id   int
		want bool
	}{
		{
			name: "Valid id",
			id:   1,
			want: true,
		},
		{
			name: "Non-existent id",
			id:   2,
			want: false,
		},
		{
			name: "Zero id",
			id:   0,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := NewTestDB(t)

			m := UserModel{db}

			exists, err := m.Exists(tt.id)

			assert.Equal(t, exists, tt.want)
			assert.NilError(t, err)
		})
	}
}
