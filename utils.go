package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"strings"

	"github.com/DavidSchott/chitchat/data"
)

// version
func version() string {
	return "0.1"
}

// Configuration stores config info of server
type Configuration struct {
	Address      string
	ReadTimeout  int64
	WriteTimeout int64
	Static       string
}

var config Configuration
var logger *log.Logger

// Convenience function for printing to stdout
func p(a ...interface{}) {
	fmt.Println(a...)
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
	config = Configuration{}
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalln("Cannot get configuration from file", err)
	}
}

// Convenience function to redirect to the error message page
func errorMessage(writer http.ResponseWriter, request *http.Request, msg string) {
	url := []string{"/err?msg=", msg}
	http.Redirect(writer, request, strings.Join(url, ""), 302)
}

// for logging
func info(args ...interface{}) {
	logger.SetPrefix("INFO ")
	logger.Println(args...)
}

func danger(args ...interface{}) {
	logger.SetPrefix("ERROR ")
	logger.Println(args...)
}

func warning(args ...interface{}) {
	logger.SetPrefix("WARNING ")
	logger.Println(args...)
}

// ReportSuccess is a helper function to return a JSON reponse indicating success
func ReportSuccess(w http.ResponseWriter, success bool, err string) {
	w.Header().Set("Content-Type", "application/json")
	if success {
		res := &data.Success{
			Sucess: success,
		}
		response, _ := json.Marshal(res)
		w.Write(response)
	} else {
		res := &data.Failure{
			Sucess: success,
			Error:  err,
		}
		response, _ := json.Marshal(res)
		w.Write(response)
	}
	return
}

func generateHTML(writer http.ResponseWriter, data interface{}, filenames ...string) {
	var files []string
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.html", file))
	}
	templates := template.Must(template.ParseFiles(files...))
	templates.ExecuteTemplate(writer, "layout", data)
}

func generateHTMLContent(writer http.ResponseWriter, data interface{}, file string) {
	writer.Header().Set("Content-Type", "text/html")
	t, _ := template.ParseFiles(fmt.Sprintf("templates/content/%s.html", file))
	t.Execute(writer, data)
}

// convenience function to be chained with another HandlerFunc
// that prints to the console which handler was called.
func logConsole(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
		fmt.Println("Handler function called - " + name)
		h(w, r)
	}
}

// convenience function to be chained with another HandlerFunc
// Checks if streaming via Server-Side Events is supported by the device
func checkStreamingSupport(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, ok := w.(http.Flusher)

		// Check if streaming is supported
		if !ok {
			http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		h(w, r)
	}
}

func formatEventData(msg string, user string, color string) (data []byte) {

	/*f := map[string]interface{}{
		"data": map[string]interface{}{
			"msg":   msg,
			"name":  user,
			"color": color,
		},
	}
	p(f)
	if data, err := json.Marshal(f); err != nil {
		warning("Error encoding JSON", f)
	} else {
		p(data)
		return data
	}*/
	json := strings.Join([]string{
		fmt.Sprintf("data: {\"msg\": \"%s\",", msg),
		fmt.Sprintf("\"name\": \"%s\",", user),
		fmt.Sprintf("\"color\": \"%s\"}\n", color),
		"\n\n",
	}, "")
	p(json)
	data = []byte(json)
	/*json := strings.Join([]string{
		"data: {\n",
		"data: \"msg\": \"" + msg + "\",\n",
		"data: \"name\": \"" + user + "\",\n",
		"data: \"color\": \"" + color + "\"\n",
		"data: }\n\n",
	}, "")
	*/
	return
}
