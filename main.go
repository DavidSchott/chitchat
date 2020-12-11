package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/DavidSchott/chitchat/handler"
)

func main() {
	//handler.Init()
	// handle static assets by routing requests from /static/ => "public" directory
	staticDir := "/static/"
	handler.Mux.PathPrefix(staticDir).Handler(http.StripPrefix(staticDir, http.FileServer(http.Dir(handler.Config.Static))))

	// Determine address to listen on
	var address string
	port := os.Getenv("PORT")
	// If port is not set, Heroku is not being used
	if port == "" {
		port = "80"
		address = handler.Config.Address
	} else {
		// If port is set, Heroku is being used
		address = "0.0.0.0:" + port
	}

	// starting up the server
	server := &http.Server{
		Addr:           address,
		Handler:        handler.Mux,
		ReadTimeout:    time.Duration(handler.Config.ReadTimeout * int64(time.Second)),
		WriteTimeout:   time.Duration(handler.Config.WriteTimeout * int64(time.Second)),
		MaxHeaderBytes: 1 << 20,
	}
	fmt.Println("ChitChat", version(), "started at", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Error starting server", err.Error())
	}
	/* TLS is already enabled on PaaS platform, so this is commented out:
	if err := server.ListenAndServeTLS("gencert/cert.pem", "gencert/key.pem"); err != nil {
		fmt.Println("Error starting server", err.Error())
	}*/
}

// version
func version() string {
	return "0.3"
}
