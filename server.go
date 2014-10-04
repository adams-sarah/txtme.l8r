package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"
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
	ok := validateReqMethod(w, req)
	if !ok {
		return
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(body))

	tmpDuration, _ := time.ParseDuration("2h45m")

	err = writeXMLResponse(w, tmpDuration)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func validateReqMethod(w http.ResponseWriter, req *http.Request) bool {
	if req.Method != "POST" {
		errMsg := fmt.Sprintf(
			"Invalid request method, %s. Please use POST.",
			req.Method,
		)

		http.Error(w, errMsg, http.StatusBadRequest)
		return false
	}

	return true
}

func writeXMLResponse(w http.ResponseWriter, duration time.Duration) error {
	w.Header().Set("Content-Type", "text/xml")

	tmplParams := map[string]string{"L8rDuration": duration.String()}
	tmpl, err := template.ParseFiles("templates/message_response.xml")
	if err != nil {
		return err
	}

	return tmpl.Execute(w, tmplParams)
}
