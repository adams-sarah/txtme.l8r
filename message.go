package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/sfreiberg/gotwilio"
)

const (
	MESSAGE_PART_DELAY = iota
	MESSAGE_PART_REMINDER
)

const MESSAGE_PART_DELIMITER = ":"

type Message struct {
	RawMessage url.Values
	Body       MessageBody
	UserNumber string
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
			http.StatusInternalServerError,
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
	if ew != nil {
		return
	}

	ew = (&msg).setUserNumber()

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

func (msg *Message) setUserNumber() *ErrorWriter {
	msg.UserNumber = msg.RawMessage.Get("From")
	if len(msg.UserNumber) == 0 {
		return NewErrorWriter(MissingFromErrorMsg, http.StatusBadRequest, false, nil)
	}

	return nil
}

func (msg *Message) sendLater() {
	time.Sleep(msg.Body.DelayTime)

	twilio := gotwilio.NewTwilioClient(
		os.Getenv("TWILIO_ACCOUNT_SID"),
		os.Getenv("TWILIO_AUTH_TOKEN"),
	)

	smsResponse, ex, err := twilio.SendSMS(
		os.Getenv("TWILIO_PHONE_NUMBER"),
		msg.UserNumber,
		msg.Body.ReminderMessage,
		os.Getenv("BASE_URI")+"/message/status",
		"",
	)

	if ex != nil {
		log.Printf("Error: Message.sendLater: 'twilio.Exception' returned: %#v\n", *ex)
	}

	if err != nil {
		log.Println("Error: Message.sendLater: ", err.Error())
	}

	if err == nil && ex == nil {
		log.Printf("Sent SMS. Response: %#v\n", *smsResponse)
	}
}
