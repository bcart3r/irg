package irg

import (
	"fmt"
	"github.com/bcart3r/irg/irc"
	"regexp"
)

type Bot struct {
	Nick, Name string
	Events     chan string
	Conn       *irc.Irc
}

func (b *Bot) Login(nick, name string) {
	b.setNickName(nick, name)
	b.Conn.W <- "NICK " + b.Nick
	b.Conn.W <- "NAME " + b.Nick + " 0 * :" + b.Name
}

func (b *Bot) setNickName(nick, name string) {
	b.Nick = nick
	b.Name = name
}

func Connect(server string) *Bot {
	conn := irc.Connect(server)
	events := make(chan string, 1000)

	return &Bot{"GoBot", "GoBot", events, conn}
}

func (b *Bot) RunLoop() {
	b.Conn.ReadHandler()
	b.Conn.WriteHandler()

	for {
		ln := <-b.Conn.R
		fmt.Print(ln)
		if regexp.MustCompile("^PING").Match([]byte(ln)) {
			b.Conn.W <- "PONG " + regexp.MustCompile(":.*").FindString(ln)
		}
	}
}
