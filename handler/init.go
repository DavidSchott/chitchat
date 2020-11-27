package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

// Configuration stores config info of server
type Configuration struct {
	Address      string
	ReadTimeout  int64
	WriteTimeout int64
	Static       string
}

var Config Configuration

// SetUp will register all HTTP handlers
func SetUp() *http.ServeMux {
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
	mux.HandleFunc("/chat/sse/", checkStreamingSupport(sseHandler))

	// Check password matches room
	mux.Handle("/chat/sse/login", errHandler(login))

	// Chat Sessions (Client sent events)
	mux.HandleFunc("/chat/sse/event", checkStreamingSupport(logConsole(sseActionHandler)))

	return mux
}

func init() {
	loadConfig()
	file, err := os.OpenFile("chitchat.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file", err)
	}
	logger = log.New(file, "INFO ", log.Ldate|log.Ltime|log.Lshortfile)
}

func loadConfig() {
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatalln("Cannot open config file", err)
	}
	decoder := json.NewDecoder(file)
	Config = Configuration{}
	err = decoder.Decode(&Config)
	if err != nil {
		log.Fatalln("Cannot get configuration from file", err)
	}
}
