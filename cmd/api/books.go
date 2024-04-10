package main

import (
	"errors"
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

// getBookHandler godoc
//
//	@Summary	Get a Book
//	@Tags		books
//	@Accept		json
//	@Produce	json
//	@Param		id	path		int	true	"Book ID"
//	@Success	200	{object}	models.Book
//	@Failure	400
//	@Failure	404
//	@Failure	500
//	@Router		/v1/books/{id} [get]
func (app *application) getBookHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	book, err := app.models.Books.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"book": book}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// updateBookHandler godoc
//
//	@Summary	Update a book by providing new values
//	@Tags		books
//	@Accept		json
//	@Produce	json
//	@Param		book	body		models.Book	true	"Provide Fields to change"
//	@Success	200		{object}	models.Book
//	@Failure	400
//	@Failure	404
//	@Failure	409
//	@Failure	422
//	@Failure	500
//	@Router		/v1/books/{id} [patch]
func (app *application) updateBookHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	book, err := app.models.Books.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	var input struct {
		Title  *string  `json:"title"`
		Author *string  `json:"author"`
		Year   *int32   `json:"year"`
		Pages  *int32   `json:"pages"`
		Genres []string `json:"genres"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Title != nil {
		book.Title = *input.Title
	}
	if input.Author != nil {
		book.Author = *input.Author
	}
	if input.Year != nil {
		book.Year = *input.Year
	}
	if input.Pages != nil {
		book.Pages = *input.Pages
	}
	if input.Genres != nil {
		book.Genres = *&input.Genres
	}

	v := validator.New()
	if models.ValidateBook(v, book); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Books.Update(book)
	if err != nil {
		if errors.Is(err, models.ErrEditConflict) {
			app.editConflictResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"book": book}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listBooksHandler(w http.ResponseWriter, r *http.Request) {
}
