// +build test

// above ^ is a build constriant:
//  only compile this file when flag 'test' is given

package main

import (
	"encoding/json"
	"net/url"
)

var TestMessageBase = Message{
	RawMessage: url.Values{
		"From":          []string{"+14154564567"},
		"FromCountry":   []string{"US"},
		"NumMedia":      []string{"0"},
		"ApiVersion":    []string{"2010-04-01"},
		"SmsSid":        []string{"abcd123"},
		"To":            []string{"+16501231234"},
		"ToCountry":     []string{"US"},
		"AccountSid":    []string{"abcd123"},
		"SmsMessageSid": []string{"abcd123"},
		"ToCity":        []string{"Millbrae"},
		"ToState":       []string{"CA"},
		"ToZip":         []string{"94030"},
		"FromState":     []string{"CA"},
		"FromZip":       []string{"94110"},
		"MessageSid":    []string{"abcd123"},
		"SmsStatus":     []string{"received"},
		"FromCity":      []string{"San Francisco"},
	},
}

func TestMessage(body string) Message {
	msg := TestMessageBase
	msg.RawMessage.Set("Body", body)
	return msg
}

func URLEncodeMessage(msg Message) []byte {
	return []byte(msg.RawMessage.Encode())
}

func JSONEncodeMessage(msg Message) ([]byte, error) {
	return json.Marshal(msg.RawMessage)
}
