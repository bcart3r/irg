package irg

import (
	"fmt"
	"regexp"
	"github.com/bcart3r/irg/irc"
)

var (
	pingMatcher = regexp.MustCompile("^PING")
	netMatcher  = regexp.MustCompile(":.*")
)

type Bot struct {
	Nick, Name string
	Events     chan string
	Chan       string
	Conn       *irc.Irc
}

func Connect(server string) *Bot {
	conn := irc.Dial(server)
	events := make(chan string, 200)

	return &Bot{"GoBot", "GoBot", events, "", conn}
}

func (b *Bot) Msg(msg string) {
	b.Conn.Write("PRIVMSG " + b.Chan + " :" + msg)
}

func (b *Bot) Join(ch string) {
	b.Conn.Write("JOIN " + ch)
	b.Chan = ch
}

func (b *Bot) Login(nick, name string) {
	b.setNickName(nick, name)
	b.Conn.Write("NICK " + b.Nick)
	b.Conn.Write("USER " + b.Nick + " 0 * :" + b.Name)
}

func (b *Bot) setNickName(nick, name string) {
	b.Nick = nick
	b.Name = name
}

func (b *Bot) RunLoop() {
	for {
		ln := <-b.Conn.R
		fmt.Print(ln)
		if pingMatcher.Match([]byte(ln)) {
			b.Conn.Write("PONG " + netMatcher.FindString(ln))
		}
	}
}
