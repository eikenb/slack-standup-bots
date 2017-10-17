package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSerial(t *testing.T) {
	up := standup{"uid", "123.456", "standup"}
	serial := up.serialize()
	assert.Equal(t, serial, "uid;123.456;standup")
	up1, err := deserialize(serial)
	assert.NoError(t, err)
	assert.Equal(t, up, up1)
}
