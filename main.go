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
	// error
	mux.HandleFunc("/err", logConsole(err))

	// header
	mux.HandleFunc("/header", logConsole(header))

	// header
	mux.HandleFunc("/process", logConsole(process))

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
