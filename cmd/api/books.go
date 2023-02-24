package main

import (
	"errors"
	"fmt"
	"libercopia/internal/data"
	"libercopia/internal/validator"
	"net/http"
	"strconv"
	"strings"
)

func (app *application) createBookHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	year, _ := strconv.Atoi(r.PostFormValue("year"))
	genres := strings.Split(strings.TrimSpace(r.PostFormValue("genres")), ",")
	price, _ := strconv.ParseFloat(r.PostFormValue("price"), 64)

	book := &data.Book{
		Title:       r.PostFormValue("title"),
		Year:        year,
		Author:      r.PostFormValue("author"),
		Genres:      genres,
		Price:       price,
		Description: r.PostFormValue("description"),
	}

	v := validator.New()

	if data.ValidateBook(v, book); !v.Valid() {
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

	http.Redirect(w, r, "/account", http.StatusSeeOther)
}

func (app *application) showBookHandler(w http.ResponseWriter, r *http.Request) {

	//id, err := app.readIDParam(r)
	//if err != nil {
	//	app.notFoundResponse(w, r)
	//	return
	//}

	//book, err := app.models.Books.GetById(id)
	//
	//if err != nil {
	//	switch {
	//	case errors.Is(err, data.ErrRecordNotFound):
	//		app.notFoundResponse(w, r)
	//	default:
	//		app.serverErrorResponse(w, r, err)
	//	}
	//	return
	//}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) updateBookHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	book, err := app.models.Books.GetById(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if r.PostFormValue("title") != "" {
		book.Title = r.PostFormValue("title")
	}
	if r.PostFormValue("year") != "" {
		book.Year, _ = strconv.Atoi(r.PostFormValue("year"))
	}
	if r.PostFormValue("author") != "" {
		book.Author = r.PostFormValue("author")
	}
	if r.PostFormValue("genres") != "" {
		slices := strings.Split(strings.TrimSpace(r.PostFormValue("genres")), ",")
		book.Genres = slices
	}
	if r.PostFormValue("prices") != "" {
		book.Price, _ = strconv.ParseFloat(r.PostFormValue("price"), 64)
	}
	if r.PostFormValue("description") != "" {
		book.Description = r.PostFormValue("description")
	}

	v := validator.New()

	if data.ValidateBook(v, book); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Books.Update(book)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Books.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) listBooksHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title  string
		Genres []string
		data.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Title = app.readString(qs, "title", "")
	input.Genres = app.readCSV(qs, "genres", []string{})

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "title", "year", "genres",
		"-id", "-title", "-year", "-genres"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	//books, metadata, err := app.models.Books.GetAll(input.Title, input.Genres, input.Filters)
	//if err != nil {
	//	app.serverErrorResponse(w, r, err)
	//	return
	//}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
