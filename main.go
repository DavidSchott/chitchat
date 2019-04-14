package main

import (
	"net/http"
	"time"
)

func main() {

	// handle static assets
	mux := http.NewServeMux()
	files := http.FileServer(http.Dir(config.Static))
	mux.Handle("/static/", http.StripPrefix("/static/", files))

	// index
	mux.HandleFunc("/", logConsole(index))

	//REST-API for chat room
	mux.HandleFunc("/chat/", logConsole(handleRoom))

	// test error
	mux.HandleFunc("/err", logConsole(err))

	// test redirect
	mux.HandleFunc("/redirect", logConsole(redirect))

	// test json
	mux.HandleFunc("/json", logConsole(jsonExample))

	// test implement
	mux.HandleFunc("/todo", logConsole(notImplemented))

	// starting up the server
	server := &http.Server{
		Addr:           config.Address,
		Handler:        mux,
		ReadTimeout:    time.Duration(config.ReadTimeout * int64(time.Second)),
		WriteTimeout:   time.Duration(config.WriteTimeout * int64(time.Second)),
		MaxHeaderBytes: 1 << 20,
	}
	p("ChitChat", version(), "started at", config.Address)
	//server.ListenAndServeTLS("gencert/cert.pem", "gencert/key.pem")
	server.ListenAndServe()
}
