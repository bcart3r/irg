package irg

import (
	"bufio"
	"log"
	"net"
)

type Conn struct {
	Reader *bufio.Reader
	Writer *bufio.Writer
	In     chan string
	Out    chan string
}

func Dial(server string) *Conn {
	conn, err := net.Dial("tcp", server)
	if err != nil {
		panic(err)
	}

	c := &Conn{
		bufio.NewReader(conn),
		bufio.NewWriter(conn),
		make(chan string, 1000),
		make(chan string, 1000),
	}

	go readHandler(c)
	go writeHandler(c)

	return c
}

func writeHandler(c *Conn) {
	for {
		msg := <-c.Out
		_, err := c.Writer.WriteString(msg + "\r\n")
		if err != nil {
			log.Println("Write Error: " + err.Error())
			return
		}

		c.Writer.Flush()
	}
}

func readHandler(c *Conn) {
	for {
		ln, err := c.Reader.ReadString(byte('\n'))
		if err != nil {
			log.Println("Read Error: " + err.Error())
			return
		}

		c.In <- ln
	}
}
