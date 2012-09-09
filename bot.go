package irg

import (
	"fmt"
	"regexp"
)

var (
	pingMatcher    = regexp.MustCompile("^PING")
	netMatcher     = regexp.MustCompile(":.*")
	ircLineMatcher = regexp.MustCompile(":(.+)!.+ (.+) (#.+) :(.+)")
)

type Bot struct {
	Nick    string
	Conn    *Conn
	Plugins []*Plugin
}

/*
	Connects to the given irc server returning an instance of Bot.
*/
func Connect(server string) *Bot {
	conn := Dial(server)

	return &Bot{"GoBot", conn, nil}
}

func (b *Bot) write(msg string) {
	b.Conn.Out <- msg
}

/*
	Joins the given channel.
*/
func (b *Bot) JoinChan(ch string) {
	b.write("JOIN :" + ch)
}

/*
	Logs into the server with the given nick and user.
*/
func (b *Bot) Login(nick, user string) {
	b.Nick = nick
	b.write("NICK :" + b.Nick)
	b.write("USER " + b.Nick + " * 0 :" + user)
}

/*
 Sends a PRIVMSG containing msg to the given channel.
*/
func (b *Bot) Msg(ch, msg string) {
	b.write("PRIVMSG " + ch + " :" + msg)
}

/*
	Sends a PRIVMSG containing msg to the given user
*/
func (b *Bot) Pm(user, msg string) {
	b.write("PRIVMSG " + user + " :" + msg)
}

/*
	Adds a Plugin to the Bots Plugins slice.
*/
func (b *Bot) AddPlugin(plugin *Plugin) {
	b.Plugins = append(b.Plugins, plugin)
}

//Adds a list of Plugins to the Bots Plugins slice.
func (b *Bot) AddPlugins(plugins []*Plugin) {
	for _, plugin := range plugins {
		b.Plugins = append(b.Plugins, plugin)
	}
}

/*
	Loops endlessly reading each line placed
	into the bots In channel in the order they are received
	also iterates over the bots Plugin slice for any matches from Matcher on the current line
	if the line matches the Plugins Callback is executed.
*/
func (b *Bot) RunLoop() {
	irc := make(map[string]string)

	for {
		ln := <-b.Conn.In
		fmt.Print(ln)
		if pingMatcher.Match([]byte(ln)) {
			b.write("PONG " + netMatcher.FindString(ln))
		}

		for _, plugin := range b.Plugins {
			if plugin.Matcher.Match([]byte(ln)) {
				go func() {
					matcher := ircLineMatcher.FindStringSubmatch(ln)
					irc["user"] = matcher[1]
					irc["msgType"] = matcher[2]
					irc["chan"] = matcher[3]
					irc["msg"] = matcher[4]

					plugin.Callback(b, irc)
				}()
			}
		}
	}
}
