package main

import (
	"fmt"
	"strings"

	"github.com/nlopes/slack"
)

// Fake/Mock slack setup
var testmess *fakeMessenger

func init() {
	api = fakeSlack{}
	testmess = newFakeMessenger()
	mess = testmess
}

var (
	botuser = &slack.User{
		ID:       botid,
		Name:     botname,
		RealName: botrealname,
	}
	botuser_details = &slack.UserDetails{
		ID:   botid,
		Name: botname,
	}

	testuser, testuser_details = makeTestUser("testuserid")

	testchannel  = makeTestChan(true)
	otherchannel = makeTestChan(false)
	privchannel  = &slack.Group{}
)

func makeTestChan(member bool) *slack.Channel {
	return &slack.Channel{IsMember: false}
}

func makeTestUser(id string) (*slack.User, *slack.UserDetails) {
	name := strings.TrimSuffix(id, "id")
	user := &slack.User{
		ID:       id,
		Name:     name,
		RealName: "real " + name,
	}
	user_details := &slack.UserDetails{
		ID:   id,
		Name: name,
	}
	return user, user_details
}

type fakeSlack struct{}

func (s fakeSlack) GetChannels(bool) ([]slack.Channel, error) {
	return []slack.Channel{*testchannel}, nil
}

func (s fakeSlack) GetGroupInfo(name string) (*slack.Group, error) {
	switch name {
	case "privatechan":
		return privchannel, nil
	default:
		return nil, fmt.Errorf("not_a_channel")
	}
}

func (s fakeSlack) GetChannelInfo(name string) (*slack.Channel, error) {
	switch name {
	case botchannel:
		return nil, fmt.Errorf("not_a_channel")
	case "testchannel":
		return testchannel, nil
	default:
		panic(name)
	}
}

func (s fakeSlack) GetUserInfo(userid string) (*slack.User, error) {
	user, _ := makeTestUser(userid)
	return user, nil
}

type fakeMessenger struct {
	rtm chan slack.RTMEvent
}

func newFakeMessenger() *fakeMessenger {
	return &fakeMessenger{rtm: make(chan slack.RTMEvent)}
}

func (m fakeMessenger) eventChan() chan slack.RTMEvent {
	return m.rtm
}
func (m fakeMessenger) sendMessage(text, channel string) {}

//
func testMessageEvent(msgtxt, channel string) *slack.MessageEvent {
	return &slack.MessageEvent{
		Msg: slack.Msg{
			Channel:   channel,
			User:      "testuserid",
			Text:      msgtxt,
			Timestamp: "123.456",
		},
	}
}

func testRTMEvent(mtype string, msgevnt *slack.MessageEvent) slack.RTMEvent {
	return slack.RTMEvent{
		Type: mtype,
		Data: msgevnt,
	}
}
