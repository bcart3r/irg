package irg

type IrcMap map[string]string

type Plugin struct {
	Matcher  string
	Callback func(b *Bot, irc IrcMap)
}
