package main

import (
	"fmt"
	"net/http"

	"github.com/svenrisse/bookshelf/internal/models"
	"github.com/svenrisse/bookshelf/internal/validator"
)

// createBookHandler godoc
//
//	@Summary		Create a Book
//	@Description	create book with fields
//	@Tags			books
//	@Accept			json
//	@Produce		json
//	@Param			book	body		models.Book	true	"Add book"
//	@Success		201		{object}	models.Book
//	@Failure		400
//	@Failure		422
//	@Failure		500
//	@Router			/v1/books [post]
func (app *application) createBookHandler(w http.ResponseWriter, r *http.Request) {
	book := &models.Book{}
	err := app.readJSON(w, r, &book)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	if models.ValidateBook(v, book); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Books.Insert(book)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/books/%d", book.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"book": book}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
