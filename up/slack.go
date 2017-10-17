package main

import (
	"log"
	"os"

	"github.com/nlopes/slack"
)

// slack abstraction and init
var mess messengerer
var api slacker
var logg *log.Logger

func init() {
	logg := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	slack.SetLogger(logg)
}

type messengerer interface {
	eventChan() chan slack.RTMEvent
	sendMessage(text string, channel string)
}

type slacker interface {
	GetChannels(bool) ([]slack.Channel, error)
	GetChannelInfo(string) (*slack.Channel, error)
	GetUserInfo(string) (*slack.User, error)
}
