package main

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

////////////////////////////////////////////////////////////////////////
// Stand
var standCommand = command{
	function:   standFunc,
	help:       "stand TEXT: save TEXT as current standup",
	in_channel: true,
}

func newStandup(msg string, values url.Values) standup {
	chan_name := values.Get("channel_name")
	user_name := values.Get("user_name")
	return standup{Where: chan_name, Who: user_name,
		When: timestamp(), What: msg}
}

// allow overriding during testing
var timestamp = func() string {
	return time.Now().Format("2006.01.02-15:04:05")
}

func standFunc(msg string, values url.Values, db dbi) (string, error) {
	stand := newStandup(msg, values)
	err := db.putOne(stand)
	return "Standup.... recorded!", err
}

////////////////////////////////////////////////////////////////////////
// Append
var appendCommand = command{
	function:   appendFunc,
	help:       "append TEXT: append TEXT to current standup",
	in_channel: true,
}

func appendFunc(msg string, values url.Values, db dbi) (string, error) {
	what, err := showFunc("", values, db)
	if err != nil {
		return "", err
	}
	// strip off header ('user - date') line
	if prev := strings.SplitN(what, "\n", 2); len(prev) > 1 {
		what = prev[1]
	}
	stand := newStandup(what+"\n"+msg, values)
	if err := db.putOne(stand); err != nil {
		return "", err
	}
	return "Appended!", nil
}

////////////////////////////////////////////////////////////////////////
// Show
var showCommand = command{
	function: showFunc,
	help:     "show [all|USER]: show USER standup (USER defaults to you)",
	subcommands: map[string]command{
		"all": command{function: showAllFunc,
			help: "show all: show everyone's standup",
		},
	},
}

func showFunc(msg string, values url.Values, db dbi) (string, error) {
	var username string
	if fields := strings.Fields(msg); len(fields) > 0 {
		username = fields[0]
	} else {
		username = values.Get("user_name")
	}
	chan_name := values.Get("channel_name")
	stand, err := db.getOne(chan_name, username)
	return stand.String(), err
}

func showAllFunc(msg string, values url.Values, db dbi) (string, error) {
	chan_name := values.Get("channel_name")
	standups, err := db.getAll(chan_name)
	if err != nil {
		return "", err
	}
	results := make([]string, len(standups))
	for i, sup := range standups {
		results[i] = sup.String()
	}
	return strings.Join(results, "\n"), nil
}

////////////////////////////////////////////////////////////////////////
// Welcome command
const welcome = "Welcome humans!\n\n" +
	"I am pleased to inform you of a new standup command.... `/up`!\n" +
	"Please use it to record and view your standups.\n" +
	"See `/up help` for more.\nThank you and goodnight."

var welcomeCommand = command{
	function:   welcomeFunc,
	help:       "welcome: show welcome/intro message to channel",
	in_channel: true,
}

func welcomeFunc(msg string, values url.Values, db dbi) (string, error) {
	return welcome, nil
}

////////////////////////////////////////////////////////////////////////
// Secret Debugging command
var debugCommand = command{
	function: debugFunc,
	help:     "output whatever I want to check in live environment",
}

func debugFunc(msg string, values url.Values, db dbi) (string, error) {
	return fmt.Sprint(values), nil
}
