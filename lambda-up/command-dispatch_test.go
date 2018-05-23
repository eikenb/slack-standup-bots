package main

import (
	"fmt"
	"reflect"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDispatch(t *testing.T) {
	// higher level
	t.Run("lookup", testCommandLookup)
	t.Run("dispatch", testCommandDispatch)
}

// Test command lookup
type cltest struct {
	msg, rest string
	cmd       commandFunc
}

func compareFuncs(f1, f2 commandFunc) bool {
	return funcName(f1) == funcName(f2)
}
func funcName(f commandFunc) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

func testCommandLookup(t *testing.T) {
	test_lookups := []cltest{
		{msg: fmt.Sprint("stand\nmy stand"), rest: "\nmy stand",
			cmd: standFunc},
		{msg: "stand my stand", rest: "my stand", cmd: standFunc},
		{msg: "show", rest: "", cmd: showFunc},
		{msg: "show all", rest: "", cmd: showAllFunc},
		{msg: "show fred", rest: "fred", cmd: showFunc},
		{msg: "append more", rest: "more", cmd: appendFunc},
		// test a couple varieties of help
		{msg: "", rest: "", cmd: help(showCommand).function},
		{msg: "show help", rest: "", cmd: help(showCommand).function},
		{msg: "show all Help", rest: "",
			cmd: help(showCommand.subcommands["all"]).function},
	}
	for _, tl := range test_lookups {
		c, r := commandLookup(tl.msg)
		cf := c.function
		if !assert.NotNil(t, cf) {
			assert.True(t, compareFuncs(cf, tl.cmd), funcName(cf),
				funcName(tl.cmd))
		}
		assert.Equal(t, tl.rest, r)
	}
}

type cdtest struct {
	query string
	resp  string
}

const _query = "token=%s&channel_name=%s&user_name=%s&text=%s"

func query(text string) string {
	return fmt.Sprintf(_query, "111", "foo", "testuser", text)
}

func testCommandDispatch(t *testing.T) {

	test_dispatches := []cdtest{
		{query: query("stand hello world"),
			resp: channelResponse("Standup.... recorded!")},
		{query: query("show"),
			resp: ephemeralResponse("testuser - now\nhello world")},
	}
	db := fakeDb()
	for _, td := range test_dispatches {
		msg := commandDispatch(td.query, db)
		assert.Equal(t, td.resp, msg)
	}
	msg := commandDispatch(fmt.Sprintf(_query, "bad", "", "", ""), db)
	assert.Equal(t, msg, ephemeralResponse("Unauthorized"))
}
