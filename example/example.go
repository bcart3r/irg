package main

import (
	"github.com/bcart3r/irg"
	"regexp"
	"time"
)

func sayTime() *irg.Plugin {
	match := regexp.MustCompile("!time")
	run := func(b *irg.Bot, ch, sender, m string) {
		b.Msg(ch, sender+" the time is "+time.Now().String())
	}

	return &irg.Plugin{match, run}
}

func main() {
	bot := irg.Connect("irc.freenode.org:6667")

	bot.Login("Irg", "Irg")
	bot.JoinChan("#irgtalk")
	bot.AddPlugin(sayTime())
	bot.RunLoop()
}
