package main

import (
	"encoding/json"
)

type response struct {
	Response_type string   `json:"response_type"`
	Text          string   `json:"text"`
	Attachments   []string `json:"attachments,omitempty"`
}

func newResponse(rtype, text string, attachments ...string) response {
	return response{
		Response_type: rtype,
		Text:          text,
		Attachments:   attachments,
	}
}

func ephemeralResponse(text string, attachments ...string) string {
	return newResponse("ephemeral", text, attachments...).String()
}

func channelResponse(text string, attachments ...string) string {
	return newResponse("in_channel", text, attachments...).String()
}

func (r response) Bytes() []byte {
	data, err := json.Marshal(r)
	if err != nil {
		data = []byte(err.Error())
	}
	return data
}

func (r response) String() string {
	return string(r.Bytes())
}
