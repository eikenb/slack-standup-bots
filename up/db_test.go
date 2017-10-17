package main

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testups = []standup{
	standup{"testuser1id", "123.456", "what1"},
	standup{"testuser2id", "123.457", "what2"},
	standup{"testuser3id", "123.458", "what3"},
	standup{"testuser4id", "123.459", "what4"},
	standup{"testuser5id", "123.460", "what5"},
}

func TestDBPush(t *testing.T) {
	newdb()
	testup := testups[0]
	err := db.push(testup)
	assert.NoError(t, err, "push failed")
	up, err := db.recent("testuser1id")
	assert.NoError(t, err, "user not found")
	assert.Equal(t, testup, up)
	testup2 := standup{"testuser1id", "321.456", "what11"}
	err = db.push(testup2)
	assert.NoError(t, err, "push2 failed")
	up, err = db.recent("testuser1id")
	assert.Equal(t, testup2, up)
}

func TestDBUsers(t *testing.T) {
	newdb()
	baseusers := make([]string, len(testups))
	for i, testup := range testups {
		err := db.push(testup)
		assert.NoError(t, err, "push failed", i)
		baseusers[i] = testup.who
	}
	dbusers, err := db.users()
	assert.NoError(t, err, "user list failed")
	assert.Len(t, dbusers, len(testups))
	sort.Strings(dbusers)
	assert.Equal(t, baseusers, dbusers)
}

func TestDBRecents(t *testing.T) {
	newdb()
	for _, testup := range testups {
		err := db.push(testup)
		assert.NoError(t, err, "push failed")
	}
	dbusers, err := db.users()
	assert.NoError(t, err, "users failed")
	assert.Len(t, dbusers, len(testups))
	sort.Strings(dbusers)
	for i, user := range dbusers {
		up, err := db.recent(user)
		assert.NoError(t, err, "recent failed", i)
		assert.Equal(t, testups[i], up)
	}
}

// Fake DB implementation
type fakedb map[string][]string

func newdb() {
	db = myDb{make(fakedb)}
}

func init() {
	newdb()
}

func (f fakedb) Do(cmd string, args ...interface{}) (interface{}, error) {
	switch cmd {
	case "LPUSH":
		key, val := args[0].(string), args[1].(string)
		f[key] = append(f[key], val)
		return len(f[key]), nil
	case "LINDEX":
		key := args[0].(string)
		idx := args[1].(int) + 1
		if val, ok := f[key]; ok {
			return val[len(val)-idx], nil
		}

	case "KEYS":
		keys := make([]interface{}, 0, len(f))
		for k := range f {
			keys = append(keys, []byte(k))
		}
		return keys, nil
	default:
		panic(cmd)
	}
	return nil, nil
}
