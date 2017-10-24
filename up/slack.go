package main

import (
	"github.com/nlopes/slack"
)

// slack abstraction and init
var mess messengerer
var api slacker

func init() {
	slack.SetLogger(logger)
}

type messengerer interface {
	eventChan() chan slack.RTMEvent
	sendMessage(text string, channel string)
}

type slacker interface {
	GetChannels(bool) ([]slack.Channel, error)
	GetChannelInfo(string) (*slack.Channel, error)
	GetGroups(bool) ([]slack.Group, error)
	GetGroupInfo(string) (*slack.Group, error)
	GetUserInfo(string) (*slack.User, error)
}
