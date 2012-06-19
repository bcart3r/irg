package irc

import (
	"net"
	"bufio"
)

type Irc struct {
	R     *bufio.Reader
	W     *bufio.Writer
	Read  chan string
	Write chan string
}

func Connect(server string) *Irc {
	conn, err := net.Dial("tcp", server)
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	in := make(chan string, 1000)
	out := make(chan string, 1000)

	return &Irc{reader, writer, in, out}
}
