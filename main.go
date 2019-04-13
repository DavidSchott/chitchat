package main

import (
	"net/http"
	"time"

	"github.com/DavidSchott/chitchat/data"
)

var chatServer data.ChatServer = data.ChatServer{
	Rooms:        make(map[int]*data.ChatRoom),
	RoomsByTitle: make(map[string]*data.ChatRoom),
	Index:        0,
}

func main() {

	// handle static assets
	mux := http.NewServeMux()
	files := http.FileServer(http.Dir(config.Static))
	mux.Handle("/static/", http.StripPrefix("/static/", files))

	// index
	mux.HandleFunc("/", logConsole(index))

	//chat
	mux.HandleFunc("/create", logConsole(create))

	// error
	mux.HandleFunc("/err", logConsole(err))

	mux.HandleFunc("/redirect", logConsole(redirect))

	mux.HandleFunc("/json", logConsole(jsonExample))

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
