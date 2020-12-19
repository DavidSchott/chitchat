package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/DavidSchott/chitchat/handler"
)

func main() {
	address := handler.Config.Address
	// If port env var is set, PaaS platform (Heroku) is being used
	if port, ok := os.LookupEnv("PORT"); ok {
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
	/* TLS is already enabled on Heroku PaaS platform, so this is commented out:
	if err := server.ListenAndServeTLS("gencert/cert.pem", "gencert/key.pem"); err != nil {
		fmt.Println("Error starting server", err.Error())
	}*/
}

// version
func version() string {
	return "0.3"
}
