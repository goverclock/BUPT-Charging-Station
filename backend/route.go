package main

import (
	"buptcs/data"
	"html/template"
	"log"
	"net/http"
)

func userlogin(writer http.ResponseWriter, request *http.Request) {

}

func index(writer http.ResponseWriter, request *http.Request) {
	// _, err := session(writer, request)
	// if err == nil { // user has logged in
	// 	http.Redirect(writer, request, "/ui", http.StatusFound)
	// } else {
	// 	http.Redirect(writer, request, "/login", http.StatusFound)
	// }
}

// show the login page
func login(writer http.ResponseWriter, request *http.Request) {
	file := "public/login.html"
	t := template.Must(template.ParseFiles(file))
	t.Execute(writer, nil)
}

// including both login and register
func authenticate(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	user_name := request.PostFormValue("user_name")
	user_passwd := request.PostFormValue("user_password")

	if _, ok := request.Form["sign_in"]; ok {
		// log.Println("登录")
	} else {
		// log.Println("注册")
		user := data.User{
			Name:    request.PostFormValue("user_name"),
			Password: request.PostFormValue("user_password"),
		}
		// log.Println(user)
		if err := user.Create(); err != nil {
			// log.Println(err, "Cannot create user")
		}
		http.Redirect(writer, request, "/login", http.StatusFound)
		return
	}

	user, err := data.UserByName(user_name)
	if err != nil {
		log.Println(err, "can't find user %v", user_name)
	}
	if user.Password == data.Encrypt(user_passwd) {
		// sess, err := user.CreateSession()
		if err != nil {
			log.Println(err, "can't create session")
		}
		cookie := http.Cookie{
			Name:     "_cookie",
			// Value:    sess.Uuid,
			HttpOnly: true,
		}
		http.SetCookie(writer, &cookie)
		http.Redirect(writer, request, "/ui", http.StatusFound)
	} else {
		// wrong password
		http.Redirect(writer, request, "/errpage", http.StatusFound)
	}
}

// func logout(writer http.ResponseWriter, request *http.Request) {
// 	cookie, err := request.Cookie("_cookie")
// 	if err != http.ErrNoCookie {
// 		// log.Println(err, "Failed to get cookie")	// just a WARNING
// 		session := data.Session{Uuid: cookie.Value}
// 		session.DeleteByUUID()
// 	}
// 	http.Redirect(writer, request, "/", http.StatusFound)
// }

// charging page
// func ui(writer http.ResponseWriter, request *http.Request) {
// 	// check if user has logged in
// 	_, err := session(writer, request)
// 	if err != nil { // user hasn't logged in
// 		http.Redirect(writer, request, "/login", http.StatusFound)
// 		return
// 	}

// 	file := "public/ui.html"
// 	t := template.Must(template.ParseFiles(file))
// 	t.Execute(writer, nil)
// }

func operation(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	log.Println(request.Form)
}

func errpage(writer http.ResponseWriter, request *http.Request) {
	file := "public/error.html"
	t := template.Must(template.ParseFiles(file))
	t.Execute(writer, nil)
}
