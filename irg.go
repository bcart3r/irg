package irg

import (
	"fmt"
	"bufio"
	"regexp"
	"net"
)

type Irc struct {
	Reader *bufio.Reader
	Writer *bufio.Writer
	R      chan string
	W      chan string
}

type Bot struct {
	Nick, Name string
	Events     chan string
	Conn       *Irc
}

func (b *Bot) Login(nick, name string) {
	b.setNickName(nick, name)
	b.Conn.W <- "NICK " + b.Nick
	b.Conn.W <- "NAME " + b.Nick + " 0 * :" + b.Name
}

func (b *Bot) setNickName(nick, name string) {
	b.Nick = nick
	b.Name = name
}

func (b *Bot) Join(ch string) {
	b.Conn.W <- "JOIN " + ch
}

func Connect(server string) *Bot {
	conn, err := net.Dial("tcp", server)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	in := make(chan string, 1000)
	out := make(chan string, 1000)
	events := make(chan string, 1000)
	irc := &Irc{reader, writer, in, out}
	bot := &Bot{"GoBot", "GoBot", events, irc}

	bot.readHandler()
	bot.writeHandler()

	return bot
}

func (b *Bot) readHandler() {
	go func() {
		for {
			ln, err := b.Conn.Reader.ReadString(byte('\n'))
			if err != nil {
				fmt.Println(err)
			}
			b.Conn.R <- ln
		}
	}()
}

func (b *Bot) writeHandler() {
	go func() {
		for {
			str := <-b.Conn.W
			_, err := b.Conn.Writer.WriteString(str + "\r\n")
			if err != nil {
				fmt.Println(err)
			}

			b.Conn.Writer.Flush()
		}
	}()
}

func (b *Bot) RunLoop() {
	for {
		ln := <-b.Conn.R
		fmt.Print(ln)
		if regexp.MustCompile("^PING").Match([]byte(ln)) {
			b.Conn.W <- "PONG " + regexp.MustCompile(":.*").FindString(ln)
		}
	}
}
