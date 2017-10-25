package main

import (
	"testing"
	"time"

	"github.com/nlopes/slack"
	"github.com/stretchr/testify/assert"
)

func TestMainExit(t *testing.T) {
	me := newBot()
	eventChan := mess.eventChan()
	go loop(me)
	eventChan <- slack.RTMEvent{Type: "exit", Data: exitLoop{}}
	// this will hang if exit isn't done
}

// test connection
func TestMainConnected(t *testing.T) {
	me := newBot()
	eventChan := mess.eventChan()
	go loop(me)
	// connected event, sending it enables the bot loop
	connected := slack.RTMEvent{Type: "connected",
		Data: &slack.ConnectedEvent{Info: &slack.Info{
			User: botuser_details,
		}}}
	assert.Equal(t, me.name, "")
	assert.Equal(t, me.id, "")
	eventChan <- connected
	time.Sleep(time.Millisecond)
	assert.Equal(t, me.name, botname)
	assert.Equal(t, me.id, "botid")
	eventChan <- slack.RTMEvent{Type: "exit", Data: exitLoop{}}
}

// test that messages go to bot
func TestMainInbox(t *testing.T) {
	me := newBot()
	eventChan := mess.eventChan()
	go loop(me)
	// Not sending connected message keeps bot listen loop from starting
	msgevent := testMessageEvent("foo", "testchannel")
	testmsg := testRTMEvent("message", msgevent)
	eventChan <- testmsg
	// all message events should get sent to bot for handling
	inboxbotmsg := <-me.inbox

	assert.Equal(t, msgevent, inboxbotmsg.ev)
	assert.Equal(t, msgevent.Channel, "testchannel")
	eventChan <- slack.RTMEvent{Type: "exit", Data: exitLoop{}}
}
