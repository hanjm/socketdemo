package main

import (
	"encoding/binary"
	"errors"
	"github.com/hanjm/log"
	"io"
	"net"
)

var ErrUnexpectedLength = errors.New("unexpected length")

// Protocol 定义消息协议
// typ 1 ping
// typ 2 pong
// typ 3 bin
// typ 4 text
// typ 5 json
// typ 6 auth
// ----------------------------------------
//  ver  | typ  | len   | body   |
//  int8 | int8 | int32 | []byte |
// ----------------------------------------
type Protocol struct {
	ver  int8   // 消息协议版本
	typ  int8   // 消息类型
	len  int32  // body长度
	body []byte // 消息体
}

func (p Protocol) WriteTo(c net.Conn) error {
	// write ver
	err := binary.Write(c, binary.BigEndian, p.ver)
	if err != nil {
		return err
	}
	// write type
	err = binary.Write(c, binary.BigEndian, p.typ)
	if err != nil {
		return err
	}
	// write length
	err = binary.Write(c, binary.BigEndian, p.len)
	if err != nil {
		return err
	}
	// data
	n, err := c.Write(p.body)
	if err != nil {
		return err
	}
	if n != int(p.len) {
		return ErrUnexpectedLength
	}
	log.Debugf("write %d bytes, data:%s", p.len, p.body)
	return err
}

func (p *Protocol) ReadFrom(c net.Conn) error {
	// read ver
	err := binary.Read(c, binary.BigEndian, &p.ver)
	if err != nil {
		return err
	}
	// read type
	err = binary.Read(c, binary.BigEndian, &p.typ)
	if err != nil {
		return err
	}
	// length
	err = binary.Read(c, binary.BigEndian, &p.len)
	if err != nil {
		return err
	}
	// read body
	var body = make([]byte, p.len)
	n, err := io.LimitReader(c, int64(p.len)).Read(body)
	if n != int(p.len) {
		return ErrUnexpectedLength
	}
	p.body = body
	log.Debugf("read %d bytes, data:%s", p.len, p.body)
	return err
}
