package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/nlopes/slack"
)

// from slack
type botmsg struct {
	ev        *slack.MessageEvent
	is_direct bool
}

// to slack
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

// send reply message to outbox
func (me bot) reply(msg botmsg, texts ...string) {
	text := strings.Join(texts, "")
	me.outbox <- replymsg{channel: msg.ev.Channel, text: text}
}

// we bot still replys even when it doesn't have anything to say
// this makes testing much easier
func (me bot) noreply(msg botmsg) {
	me.outbox <- replymsg{channel: msg.ev.Channel}
}

// send reply messages to slack
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

// listen for messages from slack
func (me bot) listen(done chan struct{}) {
	for {
		select {
		case <-done:
			return
		case msg := <-me.inbox:
			// fmt.Println("inbox", msg.ev.Text, msg.ev.User)
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
				if user, err := api.GetUserInfo(msg.ev.User); err != nil {
					me.reply(msg, "Error: ", err.Error())
				} else {
					me.reply(msg, "Hello ", user.RealName)
				}
			case "stand", "standup", "up":
				userid := msg.ev.User
				up := standup{who: userid,
					when: msg.ev.Timestamp,
					what: data}
				if err := db.push(up); err != nil {
					me.reply(msg, "Error: ", err.Error())
				} else {
					me.reply(msg, "standup recorded")
					if msg.is_direct {
						for _, room := range rooms() {
							if err := me.show(room, userid); err != nil {
								me.reply(msg, "Error: ", err.Error())
							}
						}
					}
				}
			case "show", "list", "status", "stat":
				if err := me.show(msg.ev.Channel); err != nil {
					me.reply(msg, "Error: ", err.Error())
				}
			case "help", "?":
				me.reply(msg, help(me))
			default:
				me.reply(msg, shorthelp(me))
			}
		}
	}
}

// display standup info
func (me bot) show(to string, userids ...string) error {
	var err error
	if len(userids) == 0 {
		userids, err = db.users()
		if logErr(err) {
			return err
		}
	}
	ups := make([]standup, len(userids))
	for i, uid := range userids {
		up, err := db.recent(uid)
		if logErr(err) {
			return err
		}
		up.who = uid
		ups[i] = up
	}
	for _, up := range ups {
		me.outbox <- replymsg{channel: to, text: up.String()}
	}
	return nil
}

// get room/channel ID
func rooms() []string {
	var results []string
	chns, err := api.GetChannels(false)
	logErr(err)
	for _, ch := range chns {
		results = append(results, ch.ID)
	}
	grps, err := api.GetGroups(false)
	logErr(err)
	for _, g := range grps {
		results = append(results, g.ID)
	}
	return results
}

// id comparison
func (me bot) isMe(id string) bool {
	return id == me.id
}

// @addressed to me or in a private/direct channel
func (me bot) toMe(msg botmsg) bool {
	return (strings.HasPrefix(msg.ev.Text, "<@"+me.id+">") || msg.is_direct)
}

func (me bot) String() string {
	return fmt.Sprintf("{bot: id:%s name: %s}", me.id, me.name)
}
