package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/nlopes/slack"
)

type botmsg struct {
	ev         *slack.MessageEvent
	is_private bool
}

type replymsg struct {
	channel string
	text    string
}

func (r replymsg) noreply() bool {
	return r.channel == "" && r.text == ""
}

type bot struct {
	name   string
	id     string
	inbox  chan botmsg
	outbox chan replymsg
}

func newBot() *bot {
	return &bot{inbox: make(chan botmsg), outbox: make(chan replymsg)}
}
func (me *bot) whoami(user *slack.UserDetails) {
	me.name = user.Name
	me.id = user.ID
}

func (me bot) start() chan struct{} {
	if me.id == "" {
		log.Fatal("Bot started before registered.")
	}
	done := make(chan struct{})
	go me.listen(done)
	go me.speak(done)
	return done
}

func (me bot) reply(msg botmsg, texts ...string) {
	text := strings.Join(texts, "")
	me.outbox <- replymsg{channel: msg.ev.Channel, text: text}
}

func (me bot) noreply(msg botmsg) {
	me.outbox <- replymsg{channel: msg.ev.Channel}
}

func (me bot) speak(done chan struct{}) {
	for {
		select {
		case <-done:
			return
		case msg := <-me.outbox:
			if msg.noreply() {
				continue
			}
			mess.sendMessage(msg.text, msg.channel)
		}
	}
}

func (me bot) listen(done chan struct{}) {
	for {
		select {
		case <-done:
			return
		case msg := <-me.inbox:
			// fmt.Println("inbox", msg.ev.Text)
			if !me.toMe(msg) || me.isMe(msg.ev.User) {
				// fmt.Println("continue", !me.toMe(msg), me.isMe(msg.ev.User))
				me.noreply(msg)
				continue
			}
			var cmd, data string
			text := strings.TrimPrefix(msg.ev.Text, "<@"+me.id+"> ")
			if cmdarr := strings.Fields(text); len(cmdarr) > 0 {
				cmd = cmdarr[0]
				data = strings.TrimSpace(strings.TrimPrefix(text, cmd))
			}
			switch cmd {
			case "hi", "hello":
				user, _ := api.GetUserInfo(msg.ev.User)
				me.reply(msg, "Hello ", user.RealName)
			case "standup", "stand", "up":
				// save standup
				// need user-id, datetime, text
				up := standup{who: msg.ev.User,
					when: msg.ev.Timestamp,
					what: data}
				db.push(up)
				me.reply(msg, "standup recorded")
				if msg.is_private {
					broadcast(msg.ev.User)
				}
			case "status", "stat":
				me.reply(msg, "stat")
				broadcast()
				// output all saved standups for today
			default:
				// send help
				me.reply(msg, "help")
			}
		}
	}
}

func logiferr(err error) bool {
	if err != nil {
		logg.Println(err)
		return true
	}
	return false
}

func broadcast(userids ...string) {
	var err error
	if len(userids) == 0 {
		userids, err = db.users()
		if logiferr(err) {
			return
		}
	}
	ups := make([]standup, len(userids))
	for i, uid := range userids {
		up, err := db.recent(uid)
		up.who = uid
		if logiferr(err) {
			return
		}
		ups[i] = up
	}
	chns, err := api.GetChannels(false)
	if err != nil {
		logg.Println(err)
	}
	for _, ch := range chns {
		for _, up := range ups {
			mess.sendMessage(up.String(), ch.ID)
		}
	}
}

// id comparison
func (me bot) isMe(id string) bool {
	return id == me.id
}

// @addressed to me or in a private/direct channel
func (me bot) toMe(msg botmsg) bool {
	return (strings.HasPrefix(msg.ev.Text, "<@"+me.id+">") || msg.is_private)
}

func (me bot) String() string {
	return fmt.Sprintf("{bot: id:%s name: %s}", me.id, me.name)
}
