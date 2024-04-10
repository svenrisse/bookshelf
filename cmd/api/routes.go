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

	// router.HandlerFunc(http.MethodGet, "/v1/books", app.requirePermission("books:read", app.listBookHandler))
	router.HandlerFunc(
		http.MethodPost,
		"/v1/books",
		app.requirePermission("books:write", app.createBookHandler),
	)
	// router.HandlerFunc(http.MethodGet, "/v1/books/:id", app.requirePermission("books:read", app.showBookHandler))
	// router.HandlerFunc(http.MethodPatch, "/v1/books/:id", app.requirePermission("books:write", app.updateBookHandler))
	// router.HandlerFunc(http.MethodDelete, "/v1/books/:id", app.requirePermission("books:write", app.deleteBookHandler))

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)

	router.HandlerFunc(
		http.MethodPost,
		"/v1/tokens/authentication",
		app.createAuthenticationTokenHandler,
	)

	router.Handler(http.MethodGet, "/v1/debug/vars", expvar.Handler())

	return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router)))))
}
