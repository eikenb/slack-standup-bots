package main

import (
	"fmt"

	"github.com/nlopes/slack"
)

// handle RTM event stream in main loop
// forward on messages to bot
func main() {
	me := newBot()
	loop(me)
}

// used with testing to exit main loop
type exitLoop struct{}

func loop(me *bot) {
	for msg := range mess.eventChan() {
		switch ev := msg.Data.(type) {
		case *slack.ConnectedEvent:
			me.whoami(ev.Info.User)
			me.start()
		case *slack.MessageEvent:
			ch, _ := api.GetChannelInfo(ev.Channel)
			is_private := (ch == nil) // private/direct chans have no info
			// fmt.Printf("Message: %v\n", ev)
			me.inbox <- botmsg{ev, is_private}
		case exitLoop:
			return // Used in testing
		default:
			// XXX log this?
			fmt.Printf("Event: %T; %v\n", ev, ev)
		}
	}
}
