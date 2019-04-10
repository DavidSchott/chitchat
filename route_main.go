package main

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"
)

// GET /err?msg=
// shows the error message page
func err(writer http.ResponseWriter, request *http.Request) {
	vals := request.URL.Query()
	fmt.Fprintf(writer, "Error: %s!", vals.Get("msg"))
	warning(fmt.Sprintf("Error: %s!", vals.Get("msg")))
}

func index(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "root url")
	info("accessed ")
}

// GET /header
// shows headers
func header(w http.ResponseWriter, r *http.Request) {
	for k, v := range r.Header {
		fmt.Fprintln(w, fmt.Sprintf("%s = %s", k, v))
	}
}

// POST /process
// prints posted form
func process(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(1024)
	fmt.Fprintln(w, r.MultipartForm)
	fmt.Fprintln(w, r.MultipartForm.Value["number"])
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
