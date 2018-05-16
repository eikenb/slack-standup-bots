package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResponse(t *testing.T) {
	r := channelResponse("foo")
	assert.Equal(t, `{"response_type":"in_channel","text":"foo"}`, r)
	r = ephemeralResponse("foo")
	assert.Equal(t, `{"response_type":"ephemeral","text":"foo"}`, r)
}
