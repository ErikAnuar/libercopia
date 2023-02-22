package main

import (
	"html/template"
	"net/http"
)

type ViewData struct {
	Title   string
	Message string
}

func (app *application) openPage(w http.ResponseWriter, r *http.Request) {
	data := ViewData{
		Title:   "",
		Message: "",
	}
	tmpl, _ := template.ParseFiles("ui/templates/page.html")
	tmpl.Execute(w, data)
}