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
			logger.Printf("Connected, Initializing bot...")
			me.whoami(ev.Info.User)
			me.start()
			logger.Printf("Ready\n")
		case *slack.MessageEvent:
			// fmt.Printf("Message: %v\n", ev)
			me.inbox <- botmsg{ev, isDirectMessage(ev.Channel)}
		case exitLoop:
			return // Used in testing
		case *slack.LatencyReport:
			logger.Printf("Event: %T; %v\n", ev, ev)
		default: // ignore everything else
		}
	}
}

// if channel has channelinfo, it is public
// if channel has group info, it is private or group DM channel
// if both are nil, it is a direct message channel
func isDirectMessage(channel string) bool {
	ch, err := api.GetChannelInfo(channel)
	if err != nil && err.Error() != "channel_not_found" {
		logger.Println(err)
	}
	gr, err := api.GetGroupInfo(channel)
	if err != nil && err.Error() != "channel_not_found" {
		logger.Println(err)
	}
	return (ch == nil && gr == nil)
}

func healthcheck() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprintf(w, "ok"); err != nil {
			logErr(err)
		}
	})
	logger.Printf("Healthcheck endpoint up.")
	fatalErr(http.ListenAndServe(":8080", nil))
}
