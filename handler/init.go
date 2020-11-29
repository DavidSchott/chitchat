package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/DavidSchott/chitchat/data"
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
var Mux *http.ServeMux

// registerHandlers will register all HTTP handlers
func registerHandlers() *http.ServeMux {
	mux := http.NewServeMux()
	// index
	mux.HandleFunc("/", logConsole(index))
	//"about" page
	mux.Handle("/about", logConsole(about))

	// Random junk for experimentation
	//mux.Handle("/test", errHandler(test))
	// test error
	//mux.HandleFunc("/err", logConsole(err))

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
	mux.HandleFunc("/chat/sse/", checkStreamingSupport(errHandler(sseHandler)))

	// Check password matches room
	mux.Handle("/chat/sse/login", errHandler(login))

	// Chat Sessions (Client sent events)
	mux.HandleFunc("/chat/sse/event", checkStreamingSupport(errHandler(sseActionHandler)))

	return mux
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
