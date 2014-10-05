package main

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	MESSAGE_PART_DELAY = iota
	MESSAGE_PART_REMINDER
)

const MESSAGE_PART_DELIMITER = ":"

type Message struct {
	RawMessage url.Values
	Body       MessageBody
}

type MessageBody struct {
	ReminderMessage string
	DelayTime       time.Duration
}

func decodeMessage(w http.ResponseWriter, req *http.Request) (msg Message, ew *ErrorWriter) {

	b, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()

	if err != nil {
		ew = NewErrorWriter(
			ParseRequestBodyErrorMsg,
			http.StatusBadRequest,
			false,
			err,
		)

		return
	}

	msg.RawMessage, err = url.ParseQuery(string(b))
	if err != nil {
		ew = NewErrorWriter(
			ParseRequestBodyErrorMsg,
			http.StatusBadRequest,
			false,
			err,
		)

		return
	}

	ew = (&msg).decodeBody()

	return
}

func (msg *Message) decodeBody() *ErrorWriter {
	msgParts := strings.Split(msg.RawMessage.Get("Body"), MESSAGE_PART_DELIMITER)

	if len(msgParts) < 2 {
		return NewErrorWriter(DefaultMessageBodyErrorMsg, http.StatusBadRequest, true, nil)
	}

	var err error
	msg.Body.DelayTime, err = time.ParseDuration(msgParts[MESSAGE_PART_DELAY])
	if err != nil {
		return NewErrorWriter(ParseDurationErrorMsg, http.StatusBadRequest, true, err)
	}

	// Add back any semicolons we split on in the reminder message
	msg.Body.ReminderMessage = strings.Join(
		msgParts[MESSAGE_PART_REMINDER:],
		":",
	)

	return nil
}
