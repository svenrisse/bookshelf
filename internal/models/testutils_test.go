package models

import (
	"database/sql"
	"os"
	"testing"
)

func NewTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open(
		"postgres",
		"postgres://testuser:pa55word@127.0.0.1/test_bookshelf?sslmode=disable",
	)
	if err != nil {
		t.Fatal(err)
	}

	script, err := os.ReadFile("./testdata/setup.sql")
	if err != nil {
		db.Close()
		t.Fatal(err)
	}
	_, err = db.Exec(string(script))
	if err != nil {
		db.Close()
		t.Fatal(err)
	}

	t.Cleanup(func() {
		defer db.Close()

		script, err := os.ReadFile("./testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}
		_, err = db.Exec(string(script))
		if err != nil {
			t.Fatal(err)
		}
	})

	return db
}
