package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewErrorWriterNoUnderlyingError(t *testing.T) {
	ew := NewErrorWriter("Game over.", http.StatusBadRequest, true, nil)

	if ew.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected new ErrorWriter to have a Code == %d. Got %d", http.StatusBadRequest, ew.StatusCode)
	}

	if ew.ClientMessage != "Game over." {
		t.Fatalf("Expected new ErrorWriter to have a ClientMessage == %s. Got %s", "Game over.", ew.ClientMessage)
	}

	if !ew.IsUserFacing {
		t.Fatalf("Expected new ErrorWriter to have an IsUserFacing == %v. Got %v", true, ew.IsUserFacing)
	}

	if ew.OriginalMessage != "" {
		t.Fatalf("Expected new ErrorWriter to have an OriginalMessage == ''. Got %s", ew.OriginalMessage)
	}

}

func TestNewErrorWriterWithUnderlyingError(t *testing.T) {
	err := errors.New("This is the original error.")
	ew := NewErrorWriter("Game over.. x2.", http.StatusInternalServerError, false, err)

	if ew.StatusCode != http.StatusInternalServerError {
		t.Fatalf("Expected new ErrorWriter to have a Code == %d. Got %d", http.StatusInternalServerError, ew.StatusCode)
	}

	if ew.ClientMessage != "Game over.. x2." {
		t.Fatalf("Expected new ErrorWriter to have a ClientMessage == %s. Got %s", "Game over.. x2.", ew.ClientMessage)
	}

	if ew.IsUserFacing {
		t.Fatalf("Expected new ErrorWriter to have an IsUserFacing == %v. Got %v", false, ew.IsUserFacing)
	}

	if ew.OriginalMessage != "This is the original error." {
		t.Fatalf("Expected new ErrorWriter to have an OriginalMessage == '%s'. Got %s", "This is the original error.", ew.OriginalMessage)
	}

}

func TestErrorWriterUserFacingWriteTo(t *testing.T) {
	w := httptest.NewRecorder()

	ew := NewErrorWriter("Game over.", http.StatusBadRequest, true, nil)

	ew.WriteTo(w)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected user-facing error response to have a status code == %d. Got %d", http.StatusOK, w.Code)
	}

	if w.Header().Get("Content-Type") != "text/xml" {
		t.Fatalf("Expected user-facing error response to have a content-type == '%s'. Got '%s'", "text/xml", w.Header().Get("Content-Type"))
	}

	body := string(w.Body.Bytes())

	if !strings.HasPrefix(body, "<?xml") {
		t.Fatalf("Expected response body to be xml.")
	}

	if !strings.Contains(body, "Game over.") {
		t.Fatalf(`Expected response body to contain "Try something like,". Body was %s`, body)
	}
}

func TestErrorWriterTwilioFacing500WriteTo(t *testing.T) {
	w := httptest.NewRecorder()

	err := errors.New("This is the original error.")
	ew := NewErrorWriter("Game over. You should practice.", http.StatusInternalServerError, false, err)

	ew.WriteTo(w)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("Expected user-facing error response to have a status code == %d. Got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestErrorWriterTwilioFacing400WriteTo(t *testing.T) {
	w := httptest.NewRecorder()

	err := errors.New("This is the original error.")
	ew := NewErrorWriter("Game over. You should practice.", http.StatusBadRequest, false, err)

	ew.WriteTo(w)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected user-facing error response to have a status code == %d. Got %d", http.StatusBadRequest, w.Code)
	}
}
