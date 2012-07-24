package irg

import "regexp"

type Plugins []Plugin
type Plugin struct {
	Matcher *regexp.Regexp
	Runner  func(b *Bot, ch, msg string)
}
