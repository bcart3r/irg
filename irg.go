package irg

import (
	"fmt"
	"regexp"
	"github.com/bcart3r/irg/irc"
)

var (
	pingMatcher = regexp.MustCompile("^PING")
	netMatcher  = regexp.MustCompile(":.*")
	msgMatcher  = regexp.MustCompile(" :.+")
	chanMatcher = regexp.MustCompile("#\\w+")
)

type Bot struct {
	Nick, Name string
	Plugins    []*Plugin
	Channels   []string
	Conn       *irc.Irc
}

type Plugin struct {
	Matcher *regexp.Regexp
	Run     func(msg string) string
}

/*
Connects to the given server returning
a Pointer to a Bot struct.
*/
func Connect(server string) *Bot {
	conn := irc.Dial(server)

	return &Bot{"GoBot", "GoBot", nil, "", conn}
}

/*
Writes a PRIVMSG to the Bots current irc channel.
*/
func (b *Bot) Msg(ln, msg string) {
	ch := chanMatcher.FindString(ln)
	b.Conn.Write("PRIVMSG " + ch + " :" + msg)
}

/*
Pm's the given user the supplied msg.
*/
func (b *Bot) Pm(user, msg string) {
	b.Conn.Write("PRIVMSG " + user + " :" + msg)
}

/*
Joins the given irc channel then
sets the Bot's Chan var to the
value of ch.
*/
func (b *Bot) Join(ch string) {
	b.Conn.Write("JOIN " + ch)
	b.Channels = append(b.Channels, ch)
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
Appends a Plugin struct to the bots Plugin slice.
*/
func (b *Bot) AddPlugin(plugin *Plugin) {
	b.Plugins = append(b.Plugins, plugin)
}

/*
Main loop that blocks Reading each line given
from the server reacting on that line if neccesary.
*/
func (b *Bot) RunLoop() {
	for {
		ln := <-b.Conn.R
		fmt.Print(ln)
		if pingMatcher.Match([]byte(ln)) {
			b.Conn.Write("PONG " + netMatcher.FindString(ln))
		}

		for _, plugin := range b.Plugins {
			go func() {
				if plugin.Matcher.Match([]byte(ln)) {
					b.Msg(ln, plugin.Run(msgMatcher.FindString(ln)))
				}
			}()
		}
	}
}
