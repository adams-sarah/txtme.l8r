package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Success
func TestMessageHandler200(t *testing.T) {
	var buf bytes.Buffer
	buf.Write(URLEncodeMessage(TestMessage("20h: eat foobar")))

	request, err := http.NewRequest("POST", "/message.xml", &buf)
	if err != nil {
		t.Fatalf("Expected http.NewRequest not to return an error. Got: %v", err)
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	MessageHandler(w, request)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected response code to be %d. Got %d", http.StatusOK, w.Code)
	}

	respContentType := w.Header().Get("Content-Type")

	if respContentType != "text/xml" {
		t.Fatalf("Expected content-type header to be 'text/xml'. Got %s", respContentType)
	}

	respBody := string(w.Body.Bytes())

	if !strings.HasPrefix(respBody, "<?xml") {
		t.Fatalf("Expected response body to be xml.")
	}

	if !strings.Contains(respBody, "We'll text you in 20h") {
		t.Fatalf(`Expected response body to contain "We'll text you in 20h". Body was %s`, respBody)
	}

}

// Twilio sends a GET request (instead of POST)
// We should respond with a 400
func TestMessageHandlerGET(t *testing.T) {
	var buf bytes.Buffer
	buf.Write(URLEncodeMessage(TestMessage("20h: eat foobar")))

	request, err := http.NewRequest("GET", "/message.xml", &buf)
	if err != nil {
		t.Fatalf("Expected http.NewRequest not to return an error. Got: %v", err)
	}

	w := httptest.NewRecorder()
	MessageHandler(w, request)

	if w.Code != http.StatusBadRequest {
		t.Fatalf(
			"Expected response code to be %d. Got %d",
			http.StatusBadRequest,
			w.Code,
		)
	}
}

// Twilio sends JSON data
// We should respond with a 400
func TestMessageHandlerJSONBody(t *testing.T) {
	jsonBody, err := JSONEncodeMessage(TestMessage("20h: eat foobar"))
	if err != nil {
		t.Fatalf("Expected JSONEncodeMessage not to return an error. Got: %v", err)
	}

	var buf bytes.Buffer
	buf.Write(jsonBody)

	request, err := http.NewRequest("POST", "/message.xml", &buf)
	if err != nil {
		t.Fatalf("Expected http.NewRequest not to return an error. Got: %v", err)
	}

	request.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	MessageHandler(w, request)

	if w.Code != http.StatusBadRequest {
		t.Fatalf(
			"Expected response code to be %d. Got %d",
			http.StatusBadRequest,
			w.Code,
		)
	}
}

// User sends an SMS message without a MESSAGE_PART_DELIMITER
// We should respond with a 200,
// and the body should be an xml error message for the user
func TestMessageHandlerMalformedSMS(t *testing.T) {
	var buf bytes.Buffer
	buf.Write(URLEncodeMessage(TestMessage("eat foobar")))

	request, err := http.NewRequest("POST", "/message.xml", &buf)
	if err != nil {
		t.Fatalf("Expected http.NewRequest not to return an error. Got: %v", err)
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	MessageHandler(w, request)

	if w.Code != http.StatusOK {
		t.Fatalf(
			"Expected response code to be %d. Got %d",
			http.StatusOK,
			w.Code,
		)
	}

	respContentType := w.Header().Get("Content-Type")

	if respContentType != "text/xml" {
		t.Fatalf("Expected content-type header to be 'text/xml'. Got %s", respContentType)
	}

	respBody := string(w.Body.Bytes())

	if !strings.HasPrefix(respBody, "<?xml") {
		t.Fatalf("Expected response body to be xml.")
	}

	if !strings.Contains(respBody, "Try something like,") {
		t.Fatalf(`Expected response body to contain "Try something like,". Body was %s`, respBody)
	}
}
