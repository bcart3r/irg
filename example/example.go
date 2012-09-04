package main

import (
	"irg"
	"regexp"
	"time"
)

func sayTime() *irg.Plugin {
	match := regexp.MustCompile("!time")
	run := func(b *irg.Bot, irc map[string]string) {
		b.Msg(irc["chan"],
			"hey there "+irc["sender"]+" the time is "+time.Now().String())
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
