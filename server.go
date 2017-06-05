package main

import (
	"github.com/hanjm/log"
	"io"
	"net"
)

func handler(c net.Conn) {
	defer c.Close()
	log.Infof("new conn from %s", c.RemoteAddr())
	for {
		msg := new(Protocol)
		err := msg.ReadFrom(c)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Error("msg.ReadFrom error", err)
		}
		log.Infof("received %d bytes, type:%d body:%v %s", msg.len, msg.typ, msg.body, msg.body)
	}
	log.Infof("conn closed")
}

func main() {
	l, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Error("net.Listen error")
		return
	}
	log.Info("tcp server start at :8000")
	for {
		c, err := l.Accept()
		if err != nil {
			log.Error("accept error", err)
			break
		}
		go handler(c)
	}
}
