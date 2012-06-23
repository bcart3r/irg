package irc

import (
	"bufio"
	"net"
)

type Irc struct {
	Reader *bufio.Reader
	Writer *bufio.Writer
	R      chan string
	W      chan string
}

func Dial(server string) *Bot {
	conn, err := net.Dial("tcp", server)
	if err != nil {
		log.Panic(err)
		return
	}

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	in := make(chan string, 1000)
	out := make(chan string, 1000)
	irc := &Irc{reader, writer, in, out}

	go readHandler(irc)
	go writeHandler(irc)

	return irc
}

func (i *Irc) Write(msg string) {
	i.W <- msg
}

func writeHandler(i *Irc) {
	for {
		str := <-i.W
		_, err := i.Writer.WriteString(str + "\r\n")
		if err != nil {
			log.Println("Write Error: " + err)
			return
		}

		i.Writer.Flush()
	}
}

func readHandler(i *Irc) {
	for {
		ln, err := i.Reader.ReadString(byte('\n'))
		if err != nil {
			log.Println("Read Error: " + err)
			return
		}
		i.R <- ln
	}
}
