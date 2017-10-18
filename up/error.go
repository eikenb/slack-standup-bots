package main

import (
	"log"
	"os"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
}

func logErr(err error) bool {
	if err != nil {
		logger.Println(err)
		return true
	}
	return false
}

func fatalErr(err error) {
	if err != nil {
		logger.Fatalln(err)
	}
}
