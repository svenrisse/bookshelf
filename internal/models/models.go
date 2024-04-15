package models

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Books BookModelInterface
	Users UserModelInterface
}

func NewModels(db *sql.DB) Models {
	return Models{
		Books: BookModel{DB: db},
		Users: UserModel{DB: db},
	}
}
