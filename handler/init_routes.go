package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/DavidSchott/chitchat/data"
	"github.com/gorilla/mux"
)

// Configuration stores config info of server
type Configuration struct {
	Address      string
	ReadTimeout  int64
	WriteTimeout int64
	Static       string
}

// Config captures parsed input from config.json
var Config Configuration

// Mux contains all the HTTP handlers
var Mux *mux.Router

// registerHandlers will register all HTTP handlers
func registerHandlers() *mux.Router {
	api := mux.NewRouter()
	// TODO: Uncomment again in prod
	//api := router.Host(Config.Address).Subrouter()
	// index
	api.HandleFunc("/", logConsole(index))

	// Random junk for experimentation
	//api.Handle("/test", errHandler(test))
	// test error
	//api.HandleFunc("/err", logConsole(err))

	//REST-API for chat room [JSON]
	api.Handle("/chats", errHandler(handlePost)).Methods(http.MethodPost)
	api.Handle("/chats/{titleOrID}", errHandler(authorize(handleRoom))).Methods(http.MethodGet, http.MethodPut, http.MethodDelete)

	// List all rooms in [HTML]
	api.HandleFunc("/chats", logConsole(listChats)).Methods(http.MethodGet)

	// Entrance [HTML]
	api.Handle("/chats/{titleOrID}/entrance", errHandler(joinRoom)).Methods(http.MethodGet)

	// Check password matches room
	api.Handle("/chats/{titleOrID}/token", errHandler(login)).Methods(http.MethodPost)

	// Check password matches room
	api.Handle("/chats/{titleOrID}/token/renew", errHandler(renewToken)).Methods(http.MethodGet)

	// Load chat box
	api.HandleFunc("/chats/{titleOrID}/chatbox", logConsole(chatbox)).Methods(http.MethodGet)

	// Chat Sessions (init)
	api.HandleFunc("/chats/{titleOrID}/sse/subscribe", checkStreamingSupport(errHandler(authorize(sseHandler)))).Methods(http.MethodGet)

	// Chat Sessions (Client sent events)
	api.HandleFunc("/chats/{titleOrID}/sse/broadcast", checkStreamingSupport(errHandler(authorize(sseActionHandler)))).Methods(http.MethodPost)

	// Error page
	api.HandleFunc("/err", logConsole(handleError)).Methods(http.MethodGet)
	return api
}

func init() {
	loadConfig()
	loadLog()
	Mux = registerHandlers()
	// initialize chat server
	data.CS.Init()
}

func loadLog() {
	file, err := os.OpenFile("chitchat.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		file, err = os.OpenFile("../chitchat.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalln("Failed to open log file", err)
		}

	}
	logger = log.New(file, "INFO ", log.Ldate|log.Ltime|log.Lshortfile)
}

func loadConfig() {
	file, err := os.Open("config.json")
	if err != nil {
		file, err = os.Open("../config.json")
		if err != nil {
			log.Fatalln("Cannot open config file", err)
		}
	}
	decoder := json.NewDecoder(file)
	Config = Configuration{}
	err = decoder.Decode(&Config)
	if err != nil {
		log.Fatalln("Cannot get configuration from file", err)
	}
}
