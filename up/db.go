package main

import (
	"github.com/garyburd/redigo/redis"
)

var db myDb

// storage
type myDber interface {
	push(standup) error
	recent(who string) (standup, error)
	users() ([]string, error)
}

type redis_doer interface {
	Do(cmd string, args ...interface{}) (reply interface{}, err error)
}

type myDb struct {
	redis_doer
}

func (db myDb) push(up standup) error {
	_, err := db.Do("LPUSH", up.who, up.serialize())
	return err
}

func (db myDb) recent(who string) (standup, error) {
	val, err := redis.String(db.Do("LINDEX", who, 0))
	if err != nil {
		return standup{}, err
	}
	return deserialize(val)
}

func (db myDb) users() ([]string, error) {
	return redis.Strings(db.Do("KEYS", "*"))
}
