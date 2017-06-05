package main

import (
	"encoding/json"
	"github.com/hanjm/log"
	"net"
)

func main() {
	//conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 8000, Zone: ""})
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		log.Error("net.Dial error", err)
		return
	}
	defer conn.Close()
	msgBody, err := json.Marshal(map[string]string{"testK": "testV"})
	if err != nil {
		log.Fatal("json.Marshal error", err)
	}
	msg := Protocol{1, 5, int32(len(msgBody)), msgBody}
	msg.WriteTo(conn)
	msgBody, err = json.Marshal(map[string]string{"testK2": "testV2"})
	if err != nil {
		log.Fatal("json.Marshal error", err)
	}
	msg = Protocol{1, 5, int32(len(msgBody)), msgBody}
	msg.WriteTo(conn)
}
