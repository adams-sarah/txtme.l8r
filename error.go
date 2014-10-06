package main

import (
	"log"
	"net/http"
)

// I generally like to keep error messages+code in one place.
// This makes it easier to be consistent to the user. It also
// makes for less blood/sweat/tears when you want to change them.

type ErrorWriter struct {
	StatusCode      int
	ClientMessage   string
	OriginalMessage string
	IsUserFacing    bool
}

func (ew *ErrorWriter) Error() string {
	return ew.ClientMessage
}

const ( // HTTP request errors - only given to Twilio
	ParseRequestBodyErrorMsg   = "Unexpected request body"
	MissingFromErrorMsg        = "Missing 'From' field in request body"
	RequestMethodErrorMsg      = "Unexpected request method. Please use POST."
	RequestContentTypeErrorMsg = "Unexpected content type. Please use url-encoded form data."
)

const ServerErrorMsg = "Internal server error"

const ( // Text/SMS message body errors - given to user
	DefaultMessageBodyErrorMsg = "Hm, I didn't catch that. Try something like, '2h: Buy milk'."
	ParseDurationErrorMsg      = "Bugger, I couldn't read your time-delay. Try something like, '9h45m: Wake up!' or '5m: Don't forget, you are awesome.'"
)

func NewErrorWriter(msg string, code int, isForUser bool, err error) *ErrorWriter {
	var errMsg string
	if err != nil {
		errMsg = err.Error()
	}

	return &ErrorWriter{
		StatusCode:      code,
		ClientMessage:   msg,
		OriginalMessage: errMsg,
		IsUserFacing:    isForUser,
	}
}

func (ew *ErrorWriter) WriteTo(w http.ResponseWriter) {
	log.Printf("Error %d: %s; Client Message: %s\n", ew.StatusCode, ew.OriginalMessage, ew.ClientMessage)

	// If the error message was intended for the user, we should
	// respond to Twilio with a 200 put the user's error message
	// in the xml response.
	if ew.IsUserFacing {
		writeXMLError(w, ew)
		return
	}

	// Otherwise (if we generated the error, or Twilio did),
	// we can just indicate an error with the status code.
	http.Error(w, ew.ClientMessage, ew.StatusCode)
}
