package irc

import (
	"net"
	"bufio"
	"fmt"
)

type Irc struct {
	Reader *bufio.Reader
	Writer *bufio.Writer
	R      chan string
	W      chan string
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

func (i *Irc) ReadHandler() {
	go func() {
		for {
			ln, err := i.Reader.ReadString(byte('\n'))
			if err != nil {
				fmt.Println(err)
			}
			i.R <- ln
		}
	}()
}

func (i *Irc) WriteHandler() {
	go func() {
		for {
			str := <-i.W
			_, err := i.Writer.WriteString(str + "\r\n")
			if err != nil {
				fmt.Println(err)
			}

			i.Writer.Flush()
		}
	}()
}
