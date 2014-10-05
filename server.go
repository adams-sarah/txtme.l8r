package main

import (
	"log"
	"net/http"
	"os"
	"text/template"
)

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	http.HandleFunc("/message.xml", MessageHandler)
	http.ListenAndServe(":"+port, nil)
}

func MessageHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		NewErrorWriter(
			RequestMethodErrorMsg,
			http.StatusBadRequest,
			false,
			nil,
		).WriteTo(w)
		return
	}

	if req.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		NewErrorWriter(
			RequestContentTypeErrorMsg,
			http.StatusBadRequest,
			false,
			nil,
		).WriteTo(w)
		return
	}

	msg, ew := decodeMessage(w, req)
	if ew != nil {
		ew.WriteTo(w)
		return
	}

	ew = writeXMLSuccess(w, msg.Body)
	if ew != nil {
		ew.WriteTo(w)
	}
}

func writeXMLSuccess(w http.ResponseWriter, body MessageBody) *ErrorWriter {
	w.Header().Set("Content-Type", "text/xml")

	tmpl, err := template.ParseFiles("templates/success.xml")
	if err != nil {
		return NewErrorWriter(
			RequestMethodErrorMsg,
			http.StatusInternalServerError,
			false,
			err,
		)
	}

	err = tmpl.Execute(w, body)
	if err != nil {
		return NewErrorWriter(
			RequestMethodErrorMsg,
			http.StatusInternalServerError,
			false,
			err,
		)
	}

	return nil
}

func writeXMLError(w http.ResponseWriter, ew *ErrorWriter) {
	w.Header().Set("Content-Type", "text/xml")
	log.Println("Writing XML Error")

	tmpl, err := template.ParseFiles("templates/error.xml")
	if err != nil {
		http.Error(w, ServerErrorMsg, http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, ew)
	if err != nil {
		http.Error(w, ServerErrorMsg, http.StatusInternalServerError)
	}

	return
}
