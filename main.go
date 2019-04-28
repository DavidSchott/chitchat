package main

import (
	"net/http"
	"time"

	"github.com/DavidSchott/chitchat/data"
)

/*
func main() {
	// initialize chat server
	data.CS.Init()
	testCreate()
	testRetrieve()
}
*/

func main() {

	// handle static assets
	mux := http.NewServeMux()
	files := http.FileServer(http.Dir(config.Static))
	mux.Handle("/static/", http.StripPrefix("/static/", files))

	// index
	mux.HandleFunc("/", logConsole(index))

	// index
	//mux.HandleFunc("/test", appHandler(test))
	mux.Handle("/test", errHandler(test))

	//REST-API for chat room
	mux.Handle("/chat/", errHandler(handleRoom))

	// List all rooms / "Join a chat room"
	mux.HandleFunc("/chat/list", logConsole(listChats))

	// Join action
	mux.Handle("/chat/join/", errHandler(joinRoom))

	// Load chat box
	mux.HandleFunc("/chat/box/", logConsole(chatbox))

	// Send action
	//	mux.HandleFunc("/chat/send/", logConsole(chatHandler))

	// Chat Sessions (init)
	mux.HandleFunc("/chat/sse/", checkStreamingSupport(sseHandler))

	// Chat Sessions (Client sent events)
	mux.HandleFunc("/chat/sse/event", checkStreamingSupport(logConsole(sseActionHandler)))

	// test error
	mux.HandleFunc("/err", logConsole(err))

	// starting up the server
	server := &http.Server{
		Addr:           config.Address,
		Handler:        mux,
		ReadTimeout:    time.Duration(config.ReadTimeout * int64(time.Second)),
		WriteTimeout:   time.Duration(config.WriteTimeout * int64(time.Second)),
		MaxHeaderBytes: 1 << 20,
	}

	// initialize chat server
	data.CS.Init()
	p("ChitChat", version(), "started at", config.Address)
	server.ListenAndServeTLS("gencert/cert.pem", "gencert/key.pem")
	//server.ListenAndServe()
}
