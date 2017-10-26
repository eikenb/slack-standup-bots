package main

import (
	"os"

	"github.com/nlopes/slack"
)

type mesenger struct {
	*slack.RTM
}

func (m mesenger) eventChan() chan slack.RTMEvent {
	return m.IncomingEvents
}

func (m mesenger) sendMessage(text, channel string) {
	m.SendMessage(m.NewOutgoingMessage(text, channel))
}

func init() {
	token := os.Getenv("SLACK_TOKEN")
	// private names for non-interface calls
	_api := slack.New(token)
	_mess := mesenger{_api.NewRTMWithOptions(
		&slack.RTMOptions{UseRTMStart: false})}
	go _mess.ManageConnection()
	// now assign to inteface types
	api = _api
	mess = _mess
}
