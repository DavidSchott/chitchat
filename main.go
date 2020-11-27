package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/DavidSchott/chitchat/data"
	"github.com/DavidSchott/chitchat/handler"
)

/*
func main() {
	// initialize chat server
	data.CS.Init()
	//testCreate()
	//testRetrieve()
}*/

func main() {
	// handle static assets
	mux := handler.SetUp()
	files := http.FileServer(http.Dir(handler.Config.Static))
	mux.Handle("/static/", http.StripPrefix("/static/", files))

	// starting up the server
	server := &http.Server{
		Addr:           handler.Config.Address,
		Handler:        mux,
		ReadTimeout:    time.Duration(handler.Config.ReadTimeout * int64(time.Second)),
		WriteTimeout:   time.Duration(handler.Config.WriteTimeout * int64(time.Second)),
		MaxHeaderBytes: 1 << 20,
	}

	// initialize chat server
	data.CS.Init()
	fmt.Println("ChitChat", version(), "started at", handler.Config.Address)
	//server.ListenAndServe()
	server.ListenAndServeTLS("gencert/cert.pem", "gencert/key.pem")
}

// version
func version() string {
	return "0.1"
}
