package main

import (
	"os"
	"testing"
)

// place to do any overrides needed for testing
func TestMain(m *testing.M) {
	timestamp = func() string { return "now" }
	go func() {
		for {
			slackToken <- "111"
		}
	}()
	os.Exit(m.Run())
}
