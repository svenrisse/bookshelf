package main

import (
	"expvar"
	"net/http"

	"github.com/julienschmidt/httprouter"
	_ "github.com/svenrisse/bookshelf/docs"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.ServeFiles("/docs/*filepath", http.Dir("docs"))

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodGet, "/v1/books", app.listBooksHandler)
	router.HandlerFunc(http.MethodPost, "/v1/books", app.createBookHandler)
	router.HandlerFunc(http.MethodGet, "/v1/books/:id", app.getBookHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/books/:id", app.updateBookHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/books/:id", app.deleteBookHandler)

	router.HandlerFunc(http.MethodGet, "/v1/auth/:provider/callback", app.AuthCallbackFunction)
	router.HandlerFunc(http.MethodGet, "/v1/auth/:provider/logout", app.AuthLogout)
	router.HandlerFunc(http.MethodGet, "/v1/auth/:provider", app.Auth)

	router.Handler(http.MethodGet, "/v1/debug/vars", expvar.Handler())

	return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit((router)))))
}
