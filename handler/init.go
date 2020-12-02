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
	//"about" page
	api.Handle("/about", logConsole(about))

	// Random junk for experimentation
	//api.Handle("/test", errHandler(test))
	// test error
	//api.HandleFunc("/err", logConsole(err))

	//REST-API for chat room [JSON]
	//mux.Handle("/chats", errHandler(handleRoom))
	api.Handle("/chats", errHandler(handlePost)).Methods(http.MethodPost)
	api.Handle("/chats/{titleOrID}", errHandler(handleRoom)).Methods(http.MethodGet, http.MethodPut, http.MethodDelete)

	// List all rooms in [HTML]
	api.HandleFunc("/chats", logConsole(listChats)).Methods(http.MethodGet)

	// Entrance [HTML]
	api.Handle("/chats/{titleOrID}/entrance", errHandler(joinRoom)).Methods(http.MethodGet)

	// Load chat box
	api.HandleFunc("/chats/{titleOrID}/chatbox", logConsole(chatbox)).Methods(http.MethodGet)

	// Send action
	//	mux.HandleFunc("/chat/send/", logConsole(chatHandler))

	// Chat Sessions (init)
	api.HandleFunc("/chats/sse/", checkStreamingSupport(errHandler(sseHandler))).Methods(http.MethodPost)

	// Check password matches room
	api.Handle("/chats/sse/login", errHandler(login)).Methods(http.MethodPost)

	// Chat Sessions (Client sent events)
	api.HandleFunc("/chats/sse/event", checkStreamingSupport(errHandler(sseActionHandler))).Methods(http.MethodPost)

	return api
}

func Init() *mux.Router {
	loadConfig()
	loadLog()
	Mux = registerHandlers()
	// initialize chat server
	data.CS.Init()
	return Mux
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
