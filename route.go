package main

import (
	"html/template"
	"net/http"
)

func index(writer http.ResponseWriter, request *http.Request) {
	_, err := session(writer, request)
	if err == nil { // user has logged in
		http.Redirect(writer, request, "/ui", http.StatusFound)
	} else {
		http.Redirect(writer, request, "/login", http.StatusFound)
	}
}

// GET /login
// Show the login page
func login(writer http.ResponseWriter, request *http.Request) {
	file := "public/login.html"
	t := template.Must(template.ParseFiles(file))
	t.Execute(writer, 0)
}

func ui(writer http.ResponseWriter, request *http.Request) {
	file := "public/ui.html"
	t := template.Must(template.ParseFiles(file))
	t.Execute(writer, 0)
}
