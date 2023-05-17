package main

import (
	"buptcs/data"
	"errors"
	"net/http"
)

type Configuration struct {
	Address      string
	ReadTimeout  int64
	WriteTimeout int64
	Static       string
}

var config Configuration

func init() {
	config.Address = "0.0.0.0:8080"
	config.ReadTimeout = 10
	config.WriteTimeout = 600
	config.Static = "public"
}

// check if the user is logged in and has a session
func session(writer http.ResponseWriter, request *http.Request) (sess data.Session, err error) {
	cookie, err := request.Cookie("_cookie")
	if err == nil {
		sess = data.Session{Uuid: cookie.Value}
		if ok, _ := sess.Check(); !ok {
			err = errors.New("invalid session")
		}
	}
	return
}
