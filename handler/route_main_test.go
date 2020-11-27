package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var writer *httptest.ResponseRecorder
var mux *http.ServeMux

func TestHandleGet(t *testing.T) {
	mux = SetUp()
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "https://127.0.0.1:8080/chat/1", nil)
	mux.ServeHTTP(writer, request)

	if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
	}
}
