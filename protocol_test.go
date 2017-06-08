package main

import (
	"github.com/hanjm/log"
	"io"
	"net"
	"testing"
)

func runClient() {
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		log.Error("net.Dial error", err)
		return
	}
	defer conn.Close()
	err = WriteJson(conn, map[string]string{"testK": "testV"})
	if err != nil {
		log.Fatal("json.Marshal error", err)
	}
	err = WriteJson(conn, map[string]string{"testK2": "testV2"})
	if err != nil {
		log.Fatal("json.Marshal error", err)
	}
}

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

func runServer(serverStartChan, serverCompleteChan chan struct{}) {
	l, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Error("net.Listen error")
		return
	}
	log.Info("tcp server start at :8000")
	serverStartChan <- struct{}{}
	c, err := l.Accept()
	if err != nil {
		log.Error("accept error", err)
		return
	}
	handler(c)
	serverCompleteChan <- struct{}{}
}

func TestWriteJson(t *testing.T) {
	// server
	serverStartChan := make(chan struct{})
	serverCompleteChan := make(chan struct{})
	go runServer(serverStartChan, serverCompleteChan)
	<-serverStartChan
	// client
	runClient()
	<-serverCompleteChan
}
