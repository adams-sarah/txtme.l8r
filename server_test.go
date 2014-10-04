package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMessageHandler(t *testing.T) {
	var buf bytes.Buffer
	buf.Write([]byte(`
		{
		  "Message": "20h10m"
		}`))

	request, err := http.NewRequest("POST", "/message.xml", &buf)
	if err != nil {
		t.Fatalf("Expected http.NewRequest to return no error. Got: %v", err)
	}

	w := httptest.NewRecorder()
	MessageHandler(w, request)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected response to be 200. Got %d", w.Code)
	}
}
