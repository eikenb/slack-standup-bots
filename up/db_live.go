// +build live

package main

import (
	"os"
	"strings"

	"github.com/garyburd/redigo/redis"
)

func init() {
	host := os.Getenv("REDIS_HOST")
	if !strings.HasSuffix(host, ":6379") {
		host = host + ":6379"
	}
	conn, err := redis.Dial("tcp", host)
	if err != nil {
	}
	db = myDb{conn}
}
