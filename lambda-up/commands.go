package main

import (
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

// allow overriding during testing
var timestamp = func() string {
	return time.Now().Format("2006.01.02-15:04:05")
}

func standFunc(msg string, values url.Values, db dbi) (string, error) {
	chan_name := values.Get("channel_name")
	user_name := values.Get("user_name")
	stand := standup{Where: chan_name, Who: user_name,
		When: timestamp(), What: msg}
	err := db.putOne(stand)
	return "standup recorded", err
}

////////////////////////////////////////////////////////////////////////
// Append
var appendCommand = command{
	function: appendFunc,
	help:     "append TEXT: append TEXT to current standup",
}

func appendFunc(msg string, values url.Values, db dbi) (string, error) {
	what, err := showFunc("", values, db)
	if err != nil {
		return "", err
	}
	// strip off header (user - date) line
	if prev := strings.SplitN(what, "\n", 2); len(prev) > 1 {
		what = prev[1]
	}
	what = what + "\n" + msg
	return standFunc(what, values, db)
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
	standup, err := db.getOne(chan_name, username)
	msg = standup.String(err)
	return standup.String(err), err
}

func showAllFunc(msg string, values url.Values, db dbi) (string, error) {
	chan_name := values.Get("channel_name")
	standups, err := db.getAll(chan_name)
	if err != nil {
		return "", err
	}
	results := make([]string, len(standups))
	for i, sup := range standups {
		results[i] = sup.String(nil)
	}
	return strings.Join(results, "\n"), nil
}
