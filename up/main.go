package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/nlopes/slack"
)

// handle RTM event stream in main loop
// forward on messages to bot
func main() {
	if len(os.Args) > 1 {
		usage()
	}
	go healthcheck()
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
			ch, err := api.GetChannelInfo(ev.Channel)
			logErr(err)
			is_private := (ch == nil) // private/direct chans have no info
			// fmt.Printf("Message: %v\n", ev)
			me.inbox <- botmsg{ev, is_private}
		case exitLoop:
			return // Used in testing
		case *slack.LatencyReport:
			logger.Printf("Event: %T; %v\n", ev, ev)
		default: // ignore everything else
		}
	}
}

func healthcheck() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprintf(w, "ok"); err != nil {
			logErr(err)
		}
	})
	fatalErr(http.ListenAndServe(":8080", nil))
}
