package main

import (
	"fmt"
	"net/http"

	"github.com/svenrisse/bookshelf/internal/models"
	"github.com/svenrisse/bookshelf/internal/validator"
)

// createBookHandler godoc
// @Summary		  Create a Book
// @Description	create book with fields
// @Tags			  books
// @Accept			json
// @Produce	  	json
// @Param		  	book	body	   models.Book  true  "Add book"
// @Success     201   {object} models.Book
// @Failure     400
// @Failure     422
// @Failure     500
// @Router			/books [post]
func (app *application) createBookHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title  string   `json:"title"`
		Author string   `json:"author"`
		Year   int32    `json:"year"`
		Pages  int32    `json:"pages"`
		Genres []string `json:"genres"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	book := &models.Book{
		Title:  input.Title,
		Author: input.Author,
		Year:   input.Year,
		Pages:  input.Pages,
		Genres: input.Genres,
	}

	v := validator.New()

	if models.ValidateBook(v, book); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/books/%d", book.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"book": book}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
