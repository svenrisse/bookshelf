package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/svenrisse/bookshelf/internal/models"
	"github.com/svenrisse/bookshelf/internal/validator"
)

func (app *application) createUsersBooksHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		BookID     int64     `json:"bookID"`
		UserID     int64     `json:"userID"`
		Read       bool      `json:"read"`
		Rating     float32   `json:"rating,omitempty"`
		ReviewBody string    `json:"reviewBody,omitempty"`
		ReadAt     time.Time `json:"readAt,omitempty"`
		ReviewedAt time.Time `json:"reviewedAt,omitempty"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	userBook := &models.UserBook{
		BookID:     input.BookID,
		UserID:     input.UserID,
		Read:       input.Read,
		Rating:     input.Rating,
		ReviewBody: input.ReviewBody,
		ReadAt:     input.ReadAt,
		ReviewedAt: input.ReviewedAt,
	}

	v := validator.New()
	if models.ValidateUserBook(v, userBook); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.UserBook.Insert(userBook)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"userBook": userBook}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteUsersBooksHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.UserBook.Delete(id)
	if err != nil {
		if errors.Is(err, models.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "UserBook successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateUsersBooksHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
	}

	userBook, err := app.models.UserBook.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	var input struct {
		Version     *int
		Read        *bool
		Rating      *float32
		reviewBody  *string
		Read_at     *time.Time
		Reviewed_at *time.Time
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// ...
}

func (app *application) listUsersBooksHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title  string
		Genres []string
		Rating float32
		Read   bool
		models.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Title = app.readString(qs, "title", "")
	input.Genres = app.readCSV(qs, "genres", []string{})
	input.Rating = app.readFloat(qs, "rating", 0, v)
	input.Read = app.readBool(qs, "read", false, v)
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafeList = []string{
		"title",
		"rating",
		"read_at",
		"reviewed_at",
		"published_at",
		"-title",
		"-rating",
		"-read_at",
		"-reviewed_at",
		"-published_at",
	}

	if models.ValidateFilters(v, input.Filters); v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	userBooks, metadata, err := app.models.UserBook.List(
		input.Title,
		input.Genres,
		input.Rating,
		input.Read,
		input.Filters,
	)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	err = app.writeJSON(
		w,
		http.StatusOK,
		envelope{"userBooks": userBooks, "metadata": metadata},
		nil,
	)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
