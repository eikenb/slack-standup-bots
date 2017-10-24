package main

import (
	"os"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
)

func init() {
	host := os.Getenv("REDIS_HOST")
	if !strings.HasSuffix(host, ":6379") {
		host = host + ":6379"
	}
	conn, err := redis.Dial("tcp", host,
		redis.DialConnectTimeout(time.Second*10))
	if err != nil {
		fatalErr(err)
	}
	db = myDb{conn}
}
