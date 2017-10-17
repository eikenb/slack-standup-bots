package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type standup struct {
	who, when, what string
}

func (s standup) String() string {
	var name string
	if user, err := api.GetUserInfo(s.who); err == nil {
		name = user.RealName
	}
	if s.what == "" {
		return fmt.Sprintf("%s hasn't submitted a standup.", name)
	}
	return fmt.Sprintf("%s [%s]\n%s", name, formatTimestamp(s.when), s.what)
}

func (s standup) serialize() string {
	return fmt.Sprintf("%s;%s;%s", s.who, s.when, s.what)
}

func deserialize(ser string) (standup, error) {
	www := strings.SplitN(ser, ";", 3)
	if len(www) < 3 {
		return standup{}, fmt.Errorf("Couldn't deserialize data: %v", ser)
	}
	who, when, what := www[0], www[1], www[2]
	return standup{who, when, what}, nil
}

func formatTimestamp(t string) string {
	ts := strings.Split(t, ".")
	if len(ts) < 1 {
		return ""
	}
	s, err := strconv.Atoi(ts[0])
	if err != nil {
		return ""
	}
	tobj := time.Unix(int64(s), 0)
	return tobj.Format("2006.01.02-15:04:05")
}
