package irg

import (
	"fmt"
	"./irc"
)

type Bot struct {
	Nick, Name string
	Events     chan string
	Conn       *irc.Irc
}

func (b *Bot) readHandler() {
	go func() {
		for {
			ln, err := b.Conn.R.ReadString(byte('\n'))
			if err != nil {
				fmt.Println(err)
			}
			b.Conn.Read <- ln
		}
	}()
}

func (b *Bot) writeHandler() {
	go func() {
		for {
			str := <-b.Conn.Write
			_, err := b.Conn.W.WriteString(str + "\r\n")
			if err != nil {
				fmt.Println(err)
			}

			b.Conn.W.Flush()
		}
	}()
}
