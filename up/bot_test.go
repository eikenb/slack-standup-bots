package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const botid = "botid"
const botname = "botname"
const botrealname = "botrealname"
const botchannel = "botchannel"

var message_tests = []test_pair{
	// no response
	newTestPair("hi", "testchannel", ""),
	// directed hi
	newTestPair("<@"+botid+"> hi", "testchannel", "Hello real testuser"),
	// private channel hi
	newTestPair("hi", botchannel, "Hello real testuser"),
	// standup report
	newTestPair("<@"+botid+"> standup foo", "testchannel", "standup recorded"),
	// standup stat
}

func TestBotResponse(t *testing.T) {
	me := newBot()
	me.whoami(botuser_details)
	done := make(chan struct{})
	go me.listen(done)
	for _, tp := range message_tests {
		me.inbox <- tp.in
		reply := <-me.outbox
		assert.Equal(t, tp.out, reply)
	}
	done <- struct{}{}
}

func TestBotSave(t *testing.T) {
	me := newBot()
	me.whoami(botuser_details)
	done := me.start()
	msg := botmsg{
		ev: testMessageEvent("<@"+botid+"> standup foo", "testchannel")}
	me.inbox <- msg
	dbup, _ := db.recent("testuserid")
	testup, _ := deserialize("testuserid;123.456;foo")
	assert.Equal(t, testup, dbup)
	done <- struct{}{}
}

// helper code
type test_pair struct {
	in  botmsg
	out replymsg
}

func is_private(channel string) bool {
	return channel == botchannel
}

func newTestPair(msgtxt, channel, reptxt string) test_pair {
	return test_pair{
		in:  botmsg{testMessageEvent(msgtxt, channel), is_private(channel)},
		out: replymsg{channel, reptxt},
	}
}
