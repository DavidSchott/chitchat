package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/DavidSchott/chitchat/handler"
)

func main() {
	//handler.Init()
	// handle static assets by routing requests from /static/ => "public" directory
	staticDir := "/static/"
	handler.Mux.PathPrefix(staticDir).Handler(http.StripPrefix(staticDir, http.FileServer(http.Dir(handler.Config.Static))))

	// starting up the server
	server := &http.Server{
		Addr:           handler.Config.Address,
		Handler:        handler.Mux,
		ReadTimeout:    time.Duration(handler.Config.ReadTimeout * int64(time.Second)),
		WriteTimeout:   time.Duration(handler.Config.WriteTimeout * int64(time.Second)),
		MaxHeaderBytes: 1 << 20,
	}
	fmt.Println("ChitChat", version(), "started at", handler.Config.Address)
	//server.ListenAndServe()
	if err := server.ListenAndServeTLS("gencert/cert.pem", "gencert/key.pem"); err != nil {
		fmt.Println("Error starting server", err.Error())
	}
}

// version
func version() string {
	return "0.3"
}
