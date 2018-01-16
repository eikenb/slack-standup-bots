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

func TestBot(t *testing.T) {
	me := newBot()
	me.whoami(botuser_details)
	go me.listen()
	t.Run("response", testBotResponse(me))
	t.Run("save", testBotSave(me))
	db = myDb{make(fakedb)}
	t.Run("append", testBotAppend(me))
	db = myDb{make(fakedb)}
	t.Run("show", testBotShow(me))
	me.stop()
}

func testBotResponse(me *bot) func(*testing.T) {
	return func(t *testing.T) {
		for _, tp := range message_tests {
			me.inbox <- tp.in
			reply := <-me.outbox
			assert.Equal(t, tp.out, reply)
		}
	}
}

func testBotSave(me *bot) func(*testing.T) {
	return func(t *testing.T) {
		msg := botmsg{
			ev: testMessageEvent("<@"+botid+"> standup foo", "testchannel")}
		me.inbox <- msg
		dbup, _ := db.recent("testuserid")
		testup, _ := deserialize("testuserid;123.456;foo")
		assert.Equal(t, testup, dbup)
		<-me.outbox
	}
}

func testBotAppend(me *bot) func(*testing.T) {
	return func(t *testing.T) {
		msg := botmsg{
			ev: testMessageEvent("<@"+botid+"> standup foo", "testchannel")}
		me.inbox <- msg
		<-me.outbox
		msg = botmsg{
			ev: testMessageEvent("<@"+botid+"> append bar", "testchannel")}
		me.inbox <- msg
		dbup, _ := db.recent("testuserid")
		testup, _ := deserialize("testuserid;123.456;foo bar")
		assert.Equal(t, testup, dbup)
		<-me.outbox
	}
}

func testBotShow(me *bot) func(*testing.T) {
	rep1 := replymsg{botid, "standup recorded"}
	rep2 := replymsg{testchannel.ID, "foo"}
	rep3 := replymsg{privgroup.ID, "foo"}
	return func(t *testing.T) {
		testBotSave(me)(t)
		msg := botmsg{
			ev: testMessageEvent("standup foo", botid), is_direct: true}
		me.inbox <- msg
		reply := <-me.outbox
		assert.Equal(t, rep1, reply)
		reply = <-me.outbox
		assert.Equal(t, rep2.channel, reply.channel)
		assert.Contains(t, reply.text, rep2.text)
		reply = <-me.outbox
		assert.Equal(t, rep3.channel, reply.channel)
		assert.Contains(t, reply.text, rep3.text)
	}
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
