package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {

	router := httprouter.New()
	router.ServeFiles("/static/*filepath", http.Dir("./ui/static"))

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodGet, "/", app.openIndex)
	router.HandlerFunc(http.MethodGet, "/page", app.openPage)
	router.HandlerFunc(http.MethodGet, "/cart", app.openCart)
	router.HandlerFunc(http.MethodGet, "/tokenform", app.openTokenform)
	router.HandlerFunc(http.MethodGet, "/account", app.openAccountpage)

	router.HandlerFunc(http.MethodGet, "/v2/activation_warning", app.openActivationWarning)
	router.HandlerFunc(http.MethodGet, "/v2/authentication_warning", app.openAuthenticationWarning)
	router.HandlerFunc(http.MethodGet, "/v2/permission_warning", app.openPermissionWarning)

	router.HandlerFunc(http.MethodPost, "/v1/books/create", app.requireActivatedUser(app.createBookHandler))
	router.HandlerFunc(http.MethodGet, "/v1/books/retrieve/:id", app.requireActivatedUser(app.showBookHandler))
	router.HandlerFunc(http.MethodPost, "/v1/books/update/:id", app.requireActivatedUser(app.updateBookHandler))
	router.HandlerFunc(http.MethodPost, "/v1/books/delete/:id", app.requireActivatedUser(app.deleteMovieHandler))
	router.HandlerFunc(http.MethodGet, "/v1/books/list", app.requireActivatedUser(app.listBooksHandler))

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/users/activated", app.activateUserHandler)

	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	return app.recoverPanic(app.rateLimit(app.authenticate(router)))
}
