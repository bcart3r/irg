package irg

import (
	"fmt"
	"regexp"
)

var (
	pingMatcher = regexp.MustCompile("^PING")
	netMatcher  = regexp.MustCompile(":.*")
	msgMatcher  = regexp.MustCompile(" :.+")
	chanMatcher = regexp.MustCompile("#\\w+")
)

type Bot struct {
	Nick    string
	Conn    *Conn
	Plugins []*Plugin
}

func Connect(server string) *Bot {
	conn := Dial(server)

	return &Bot{"GoBot", conn, nil}
}

func (b *Bot) write(msg string) {
	b.Conn.Out <- msg
}

func (b *Bot) JoinChan(ch string) {
	b.write("JOIN :" + ch)
}

func (b *Bot) Login(nick, user string) {
	b.Nick = nick
	b.write("NICK :" + b.Nick)
	b.write("USER " + b.Nick + " * 0 :" + user)
}

func (b *Bot) Msg(ch, msg string) {
	b.write("PRIVMSG " + ch + " :" + msg)
}

func (b *Bot) AddPlugin(plugin *Plugin) {
	b.Plugins = append(b.Plugins, plugin)
}

func (b *Bot) RunLoop() {
	for {
		ln := <-b.Conn.In
		fmt.Print(ln)
		if pingMatcher.Match([]byte(ln)) {
			b.write("PONG " + netMatcher.FindString(ln))
		}

		for _, plugin := range b.Plugins {
			go func() {
				if plugin.Matcher.Match([]byte(ln)) {
					plugin.Runner(b, chanMatcher.FindString(ln), msgMatcher.FindString(ln))
				}
			}()
		}
	}
}
