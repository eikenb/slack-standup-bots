package main

import "io/ioutil"

func init() {
	logger.SetOutput(ioutil.Discard)
}
