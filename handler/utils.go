package handler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"reflect"
	"runtime"
	"strings"

	"github.com/DavidSchott/chitchat/data"
)

var logger *log.Logger

/* Convenience function for printing to stdout
func p(a ...interface{}) {
	fmt.Println(a...)
}*/

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

// ReportStatus is a helper function to return a JSON reponse indicating outcome success/failure
func ReportStatus(w http.ResponseWriter, success bool, err *data.APIError) {
	var res *data.Outcome
	w.Header().Set("Content-Type", "application/json")
	if success {
		res = &data.Outcome{
			Status: success,
		}
	} else {
		res = &data.Outcome{
			Status: success,
			Error:  err,
		}
	}
	response, _ := json.Marshal(res)
	if _, err := w.Write(response); err != nil {
		danger("Error writing", response)
	}
}

func generateHTML(writer http.ResponseWriter, data interface{}, filenames ...string) {
	var files []string
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.html", file))
	}
	templates := template.Must(template.ParseFiles(files...))
	if err := templates.ExecuteTemplate(writer, "layout", data); err != nil {
		danger("Error generating HTML template", data, err.Error())
	}
}

func generateHTMLContent(writer http.ResponseWriter, data interface{}, file string) {
	writer.Header().Set("Content-Type", "text/html")
	t, _ := template.ParseFiles(fmt.Sprintf("templates/content/%s.html", file))
	if err := t.Execute(writer, data); err != nil {
		danger("Error executing HTML template", data, err.Error())
	}
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
func checkStreamingSupport(h errHandler) http.HandlerFunc {
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
		if err := h(w, r); err != nil {
			warning("Error calling:", runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name())
		}
	}
}

func formatEventData(msg string, user string, color string) (data []byte) {
	json := strings.Join([]string{
		fmt.Sprintf("data: {\"msg\": \"%s\",", msg),
		fmt.Sprintf("\"name\": \"%s\",", user),
		fmt.Sprintf("\"color\": \"%s\"}\n", color),
		"\n\n",
	}, "")
	data = []byte(json)

	return
}
