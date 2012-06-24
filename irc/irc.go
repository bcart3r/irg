package irc

import (
	"bufio"
	"net"
	"log"
)

type Irc struct {
	Reader *bufio.Reader
	Writer *bufio.Writer
	R      chan string
	W      chan string
}

/*
Connects to the given server
return and instance of *Irc.
*/
func Dial(server string) *Irc {
	conn, err := net.Dial("tcp", server)
	if err != nil {
		log.Panic(err)
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

/*
Writes to the Irc conn.
*/
func (i *Irc) Write(msg string) {
	i.W <- msg
}

/*
Takes each line attempting to be written to the server
from the Bots W Chan and writes it to the irc Conn
in the order they are received.
*/
func writeHandler(i *Irc) {
	for {
		str := <-i.W
		_, err := i.Writer.WriteString(str + "\r\n")
		if err != nil {
			log.Println("Write Error: " + err.Error())
			return
		}

		i.Writer.Flush()
	}
}

/*
Takes each lines read from the Irc conn
and places them in the Irc conn R Chan
in the order they are received.
*/
func readHandler(i *Irc) {
	for {
		ln, err := i.Reader.ReadString(byte('\n'))
		if err != nil {
			log.Println("Read Error: " + err.Error())
			return
		}
		i.R <- ln
	}
}
