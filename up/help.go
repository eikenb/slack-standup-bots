package main

import (
	"fmt"
	"os"
)

func usage() {
	fmt.Println("Usage:", os.Args[0])
	fmt.Println("Requires 2 environment variables:")
	fmt.Println("\tSLACK_TOKEN - your slack bot api token.")
	fmt.Println("\tREDIS_HOST  - hostname:port of redis server.")
	os.Exit(1)
}

const help_tmpl = `
Usage: @%s command (or in private channel)
Commands:
    standup <what you did yesterday, doing today, etc.>
        How you enter your standup.
        Echo's it back in main channel if you do it in private channel.
        Free form. All 1 line.
    status
        Outputs everyone's most recent standup entries.
Eg.
@%s yesterday I worked on bug #1; today I worked on bug #2; no blockers
`

func help(me bot) string {
	return fmt.Sprintf(help_tmpl, me.name, me.name)
}
