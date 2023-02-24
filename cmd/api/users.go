package main

import (
	"errors"
	"fmt"
	"libercopia/internal/data"
	"libercopia/internal/validator"
	"net/http"
	"time"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		Name:      r.PostFormValue("name"),
		Surname:   r.PostFormValue("surname"),
		Email:     r.PostFormValue("email"),
		Role:      "user",
		Activated: false,
	}

	err = user.Password.Set(r.PostFormValue("password"))
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()

	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	token, err := app.models.Tokens.New(user.ID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	http.Redirect(w, r, "/tokenform", http.StatusSeeOther)

	app.background(func() {

		user_tmpl := app.models.Users.GetById(token.UserID)

		data := map[string]any{
			"activationToken": token.Plaintext,
			"userID":          token.UserID,
			"name":            user_tmpl.Name,
			"surname":         user_tmpl.Surname,
			"email":           user_tmpl.Email,
		}

		err = app.mailer.Send(user.Email, "user_welcome.tmpl", data)
		if err != nil {
			fmt.Println(err)
			app.logger.PrintError(err, nil)
		}
	})
}

func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	if data.ValidateTokenPlaintext(v, r.PostFormValue("tokenPlaintext")); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetForToken(data.ScopeActivation, r.PostFormValue("tokenPlaintext"))
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("token", "invalid or expired activation token")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	user.Activated = true

	err = app.models.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.models.Tokens.DeleteAllForUser(data.ScopeActivation, user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
