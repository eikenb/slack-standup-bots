package main

import (
	"fmt"
	"net/url"
	"strings"
)

// required parameters from slack
var required = [4]string{"token", "text", "channel_name", "user_name"}

// Top level of command tree
var commands = command{
	help: "[show|stand|append|help] (help works on subcommands)",
	subcommands: map[string]command{
		"show":   showCommand,
		"stand":  standCommand,
		"append": appendCommand,
	},
}

// look up, run command, return response
func commandDispatch(query string, db dbi) string {
	values, err := url.ParseQuery(query)
	if err != nil {
		return ephemeralResponse("ParseQuery failure")
	}

	if ok, missing := checkRequired(values); !ok {
		return ephemeralResponse(
			fmt.Sprintf("Missing required parameters: %v", missing))
	}

	message := values.Get("text")
	cmd, rest := commandLookup(message)
	if cmd.function == nil {
		cmd = help(commands)
	}
	// do this as late as possible to give the KMS call time to work
	if token := values.Get("token"); token != <-slackToken {
		return ephemeralResponse("Unauthorized")
	}

	text, err := cmd.function(rest, values, db)
	if err != nil {
		return ephemeralResponse(err.Error())
	}

	switch cmd.in_channel {
	case true:
		return channelResponse(text)
	default:
		return ephemeralResponse(text)
	}
}

func checkRequired(values url.Values) (bool, []string) {
	good := true
	missing := []string{}
	for _, k := range required {
		if _, ok := values[k]; !ok {
			missing = append(missing, k)
			good = false
		}
	}
	return good, missing
}

func help(cmd command) command {
	return command{
		function: func(_ string, values url.Values, _ dbi) (string, error) {
			name := values.Get("command")
			return fmt.Sprintf("Usage: %s %s", name, cmd.help), nil
		},
	}
}

func commandLookup(msg string) (command, string) {
	// change newlines to space-null to preserve newlines while still
	// splitting on them with Fields
	msg = strings.Replace(msg, "\n", " \000", -1)
	parts := strings.Fields(msg)
	cmd := commands
	var idx int
	for i, part := range parts {
		newcmd, ok := cmd.subCommand(part)
		if !ok {
			break
		}
		idx = i
		cmd = newcmd
	}
	if len(parts) > idx {
		idx += 1
	}
	rest := strings.Join(parts[idx:], " ")
	return cmd, strings.Replace(rest, "\000", "\n", -1)
}

// command lookup table
type commandFunc func(string, url.Values, dbi) (string, error)
type command struct {
	function    commandFunc
	help        string
	subcommands map[string]command
	in_channel  bool // response is in channel, public, if true
}

// set top level command to help output
func init() {
	commands.function = help(commands).function
}

func (cmd command) subCommand(cname string) (command, bool) {
	if newcmd, ok := cmd.subcommands[cname]; ok {
		return newcmd, true
	}
	if strings.ToLower(cname) == "help" {
		return help(cmd), true
	}
	return command{}, false
}
