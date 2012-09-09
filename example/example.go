package main

import (
	"github.com/bcart3r/irg"
	"regexp"
	"time"
)

func said() *irg.Plugin {
	match := regexp.MustCompile("!said")
	run := func(b *irg.Bot, irc irg.IrcMap) {
		b.Msg(irc["chan"], "hey "+irc["user"]+"you said "+irc["msg"])
	}

	return &irg.Plugin{match, run}
}

func sayTime() *irg.Plugin {
	match := regexp.MustCompile("!time")
	run := func(b *irg.Bot, irc irg.IrcMap) {
		b.Msg(irc["chan"],
			"hey there "+irc["user"]+" the time is "+time.Now().String())
	}

	return &irg.Plugin{match, run}
}

func main() {
	bot := irg.Connect("irc.freenode.org:6667")

	bot.Login("Irg", "Irg")
	bot.JoinChan("#irgtalk")
	bot.AddPlugin(sayTime())
	bot.AddPlugin(said())
	bot.RunLoop()
}
