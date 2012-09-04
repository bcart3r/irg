package irg

import "regexp"

type Plugin struct {
	Matcher *regexp.Regexp
	Runner  func(b *Bot, irc map[string]string)
}
