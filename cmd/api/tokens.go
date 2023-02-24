package main

import (
	"errors"
	"libercopia/internal/data"
	"libercopia/internal/validator"
	"net/http"
	"time"
)

func (app *application) createAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	data.ValidateEmail(v, r.PostFormValue("email"))
	data.ValidatePasswordPlaintext(v, r.PostFormValue("password"))

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetByEmail(r.PostFormValue("email"))
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.invalidCredentialsResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	match, err := user.Password.Matches(r.PostFormValue("password"))
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if !match {
		app.invalidCredentialsResponse(w, r)
		return
	}

	_, err = app.models.Tokens.New(user.ID, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	http.Redirect(w, r, "/account", http.StatusSeeOther)
}
