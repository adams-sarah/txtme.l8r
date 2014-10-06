package main

import (
	"io/ioutil"
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
	http.HandleFunc("/message/status", MessageStatusHandler)
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

	go msg.sendLater()

	ew = writeXMLSuccess(w, msg.Body)
	if ew != nil {
		ew.WriteTo(w)
	}
}

func MessageStatusHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		NewErrorWriter(
			RequestMethodErrorMsg,
			http.StatusBadRequest,
			false,
			nil,
		).WriteTo(w)
		return
	}

	b, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()

	if err != nil {
		ew := NewErrorWriter(
			ParseRequestBodyErrorMsg,
			http.StatusInternalServerError,
			false,
			err,
		)

		ew.WriteTo(w)
		return
	}

	log.Println("MessageStatusHandler: Received message: ", string(b))

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
