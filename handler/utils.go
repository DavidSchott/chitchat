package handler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"reflect"
	"runtime"

	"github.com/DavidSchott/chitchat/data"
)

var logger *log.Logger

/* Convenience function for printing to stdout
func p(a ...interface{}) {
	fmt.Println(a...)
}*/

// for logging
func Info(args ...interface{}) {
	logger.SetPrefix("INFO ")
	logger.Println(args...)
}

func Danger(args ...interface{}) {
	logger.SetPrefix("ERROR ")
	logger.Println(args...)
}

func Warning(args ...interface{}) {
	logger.SetPrefix("WARNING ")
	logger.Println(args...)
}

// ReportStatus is a helper function to return a JSON response indicating outcome success/failure
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
		Danger("Error writing", response)
	}
}

func generateHTML(writer http.ResponseWriter, data interface{}, filenames ...string) {
	var files []string
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.html", file))
	}
	templates := template.Must(template.ParseFiles(files...))
	if err := templates.ExecuteTemplate(writer, "layout", data); err != nil {
		Danger("Error generating HTML template", data, err.Error())
	}
}

func generateHTMLContent(writer http.ResponseWriter, data interface{}, file string) {
	writer.Header().Set("Content-Type", "text/html")
	t, _ := template.ParseFiles(fmt.Sprintf("templates/content/%s.html", file))
	if err := t.Execute(writer, data); err != nil {
		Danger("Error executing HTML template", data, err.Error())
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
