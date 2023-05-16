package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	log.Println("BUPTCS started at", config.Address)

	mux := http.NewServeMux()
	files := http.FileServer(http.Dir(config.Static))
	mux.Handle("/static/", http.StripPrefix("/static/", files))

	mux.HandleFunc("/", index)
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/authenticate", authenticate)
	mux.HandleFunc("/logout", logout)
	mux.HandleFunc("/ui", ui)
	mux.HandleFunc("/operation", operation)
	mux.HandleFunc("/operation/start_charge", operation)

	mux.HandleFunc("/errpage", errpage)

	server := &http.Server{
		Addr:           config.Address,
		Handler:        mux,
		ReadTimeout:    time.Duration(config.ReadTimeout * int64(time.Second)),
		WriteTimeout:   time.Duration(config.WriteTimeout * int64(time.Second)),
		MaxHeaderBytes: 1 << 20,
	}
	server.ListenAndServe()
}
