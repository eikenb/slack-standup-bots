package main

import (
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// sample query text
// token=nfTiev6lF9TJ6l3sJ4p6Sqd1&team_id=T8BDMTN5V&team_domain=text&channel_id=C8KR3CWQL&channel_name=foo&user_id=U8B8B8W9Z&user_name=testuser&command=%2Fup&text=show+me&response_url=https%3A%2F%2Fhooks.slack.com%2Fcommands%2FT7BDMTN5V%2F357957683877%2FnBMVF1jxOIViX3cv2uojcicl&trigger_id=357640056610.249463940166.38aec8c222a61dc59b6c824d9b9ddcb4"

func testValues(text ...string) url.Values {
	return url.Values{
		"channel_name": []string{"foo"},
		"user_name":    []string{"testuser"},
		"response_url": []string{"slack-hook-url"},
		"token":        []string{"slack-token"},
		"text":         text,
	}
}

func TestCommands(t *testing.T) {
	db := fakeDb()
	t.Run("stand", testStandFunc(db))
	db = populatedFakeDb()
	t.Run("append", testAppendFunc(db))
	db = populatedFakeDb() // append changes db
	t.Run("user", testShowFunc(db))
	t.Run("all", testShowAllFunc(db))
}

// Test command functions
func testStandFunc(db dbi) func(t *testing.T) {
	return func(t *testing.T) {
		values := testValues()
		out, err := standFunc("my standup", values, db)
		assert.NoError(t, err)
		assert.Equal(t, "standup recorded", out)
		chan_name := values.Get("channel_name")
		user_name := values.Get("user_name")
		s, err := db.getOne(chan_name, user_name)
		assert.NoError(t, err)
		assert.Equal(t, chan_name, s.Where)
		assert.Equal(t, user_name, s.Who)
		assert.Equal(t, "my standup", s.What)
	}
}

func testAppendFunc(db dbi) func(t *testing.T) {
	return func(t *testing.T) {
		values := testValues()
		out, err := appendFunc("more standup", values, db)
		assert.NoError(t, err)
		assert.Equal(t, "standup recorded", out)
		out, err = appendFunc("wait", values, db)
		assert.NoError(t, err)
		chan_name := values.Get("channel_name")
		user_name := values.Get("user_name")
		s, err := db.getOne(chan_name, user_name)
		assert.NoError(t, err)
		assert.Equal(t, "standup\nmore standup\nwait", s.What)
	}
}

// sups is defined in db_test.go
func testShowFunc(db dbi) func(t *testing.T) {
	return func(t *testing.T) {
		// USER provided
		s, err := showFunc("testuser", testValues(), db)
		assert.NoError(t, err)
		assert.Equal(t, sups[0].String(nil), s)
		// no USER provided
		s, err = showFunc("", testValues(), db)
		assert.NoError(t, err)
		assert.Equal(t, sups[0].String(nil), s)
		// bad USER
		s, err = showFunc("nobody", testValues(), db)
		assert.Error(t, err)
		assert.Equal(t, "ResourceNotFoundException: ", err.Error())
		assert.Equal(t, s, "")
	}
}
func testShowAllFunc(db dbi) func(t *testing.T) {
	results := make([]string, len(sups))
	for i, s := range sups {
		results[i] = s.String(nil)
	}
	return func(t *testing.T) {
		s, err := showAllFunc("", testValues(), db)
		assert.NoError(t, err)
		assert.Equal(t, strings.Join(results, "\n"), s)
	}
}
