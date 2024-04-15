package models

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/svenrisse/bookshelf/internal/validator"
)

var ErrDuplicateEmail = errors.New("duplicate email")

var AnonymousUser = &User{}

type User struct {
	ID        int `json:"id"`
	Provider  string
	Avatar    string
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateUser struct {
	Name string `json:"name"     example:"testuser"`
}

type UserModel struct {
	DB *sql.DB
}

type UserModelInterface interface {
	Insert(user *User) error
	Exists(id int) (bool, error)
	GetByEmail(email string) (*User, error)
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Name != "", "name", "must be provided")
	v.Check(len(user.Name) <= 500, "name", "must not be more than 500 bytes long")
}

func (m UserModel) Insert(user *User) error {
	query := `
    INSERT INTO users (id, name, avatar, provider)
    VALUES ($1, $2, $3, $4)
    RETURNING created_at`

	args := []any{user.ID, user.Name, user.Avatar, user.Provider}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (m UserModel) Exists(id int) (bool, error) {
	var exists bool

	stmt := "SELECT EXISTS(SELECT true FROM users WHERE id = $1)"

	err := m.DB.QueryRow(stmt, id).Scan(&exists)
	return exists, err
}

func (m UserModel) GetByEmail(email string) (*User, error) {
	query := `
        SELECT id, created_at, name
        FROM users
        WHERE email = $1`

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}
