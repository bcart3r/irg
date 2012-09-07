package irg

import "regexp"

type IrcMap map[string]string

type Plugin struct {
	Matcher  *regexp.Regexp
	Callback func(b *Bot, irc IrcMap)
}
