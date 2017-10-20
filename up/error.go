package main

import (
	"fmt"
	"log"
	"os"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
}

func logErr(err error) bool {
	if err != nil {
		logger.Output(2, fmt.Sprintln(err))
		return true
	}
	return false
}

func fatalErr(err error) {
	if err != nil {
		logger.Output(2, fmt.Sprintln(err))
		os.Exit(1)
	}
}
