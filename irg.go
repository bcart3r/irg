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

/*
Connects to the given server returning
an instance of the Bot struct.
*/
func Connect(server string) *Bot {
	conn := irc.Dial(server)
	events := make(chan string, 200)

	return &Bot{"GoBot", "GoBot", events, "", conn}
}

/*
Writes a PRIVMSG to the Bots current Channel.
*/
func (b *Bot) Msg(msg string) {
	b.Conn.Write("PRIVMSG " + b.Chan + " :" + msg)
}

/*
Pm's the given user the supplied msg.
*/
func (b *Bot) Pm(user, msg string) {
	b.Conn.Write("PRIVMSG " + user + " :" + msg)
}

/*
Joins the given channel then
sets the Bot's Chan var to the
value of ch.
*/
func (b *Bot) Join(ch string) {
	b.Conn.Write("JOIN " + ch)
	b.Chan = ch
}

/*
Set's the NICK and USER for the Bot
then write the values to the Irc conn.
*/
func (b *Bot) Login(nick, name string) {
	b.setNickName(nick, name)
	b.Conn.Write("NICK " + b.Nick)
	b.Conn.Write("USER " + b.Nick + " 0 * :" + b.Name)
}

func (b *Bot) setNickName(nick, name string) {
	b.Nick = nick
	b.Name = name
}

/*
Main loops that blocks Reading each line given
from the server reacting on that line if neccesary.
*/
func (b *Bot) RunLoop() {
	for {
		ln := <-b.Conn.R
		fmt.Print(ln)
		if pingMatcher.Match([]byte(ln)) {
			b.Conn.Write("PONG " + netMatcher.FindString(ln))
		}
	}
}
