package main

import (
	"buptcs/data"
	"html/template"
	"log"
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

// show the login page
func login(writer http.ResponseWriter, request *http.Request) {
	file := "public/login.html"
	t := template.Must(template.ParseFiles(file))
	t.Execute(writer, nil)
}

func authenticate(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	user_email := request.PostFormValue("user_name")
	user_passwd := request.PostFormValue("user_password")

	if _, ok := request.Form["sign_in"]; ok {
		// log.Println("登录")
	} else {
		// log.Println("注册")
		user := data.User{
			Email:    request.PostFormValue("user_name"),
			Password: request.PostFormValue("user_password"),
		}
		log.Println(user)
		if err := user.Create(); err != nil {
			log.Println(err, "Cannot create user")
		}
		http.Redirect(writer, request, "/login", http.StatusFound)
		return
	}

	user, err := data.UserByEmail(user_email)
	if err != nil {
		log.Println(err, "can't find user %v", user_email)
	}
	if user.Password == data.Encrypt(user_passwd) {
		sess, err := user.CreateSession()
		if err != nil {
			log.Println(err, "can't create session")
		}
		cookie := http.Cookie{
			Name:     "_cookie",
			Value:    sess.Uuid,
			HttpOnly: true,
		}
		http.SetCookie(writer, &cookie)
		http.Redirect(writer, request, "/ui", http.StatusFound)
	} else {
		// wrong password
		http.Redirect(writer, request, "/errpage", http.StatusFound)
	}
}

// charging page
func ui(writer http.ResponseWriter, request *http.Request) {
	file := "public/ui.html"
	t := template.Must(template.ParseFiles(file))
	t.Execute(writer, nil)
}

// TODO: a page to show an error message
func errpage(writer http.ResponseWriter, request *http.Request) {
	// file := "public/ui.html"
	// t := template.Must(template.ParseFiles(file))
	// t.Execute(writer, nil)
}
