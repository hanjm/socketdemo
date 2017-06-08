package main

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"hash/adler32"
	"io"
)

var ErrUnexpectedLength = errors.New("unexpected length")
var ErrUnexpectedCheckSum = errors.New("unexpected checkSum")

const (
	// Protocol.Ver
	Ver = 1
	// Protocol.typ
	PingMsg = 1
	PongMsg = 2
	AuthMsg = 3
	BinMsg  = 4
	TextMsg = 5
	JsonMsg = 6
	PbMsg   = 7
)

// Protocol 定义消息协议
// ----------------------------------------
//  ver  | typ  | len   | body   |
//  int8 | int8 | int32 | []byte |
// ----------------------------------------
type Protocol struct {
	ver      int8   // 消息协议版本
	typ      int8   // 消息类型
	len      int32  // body长度
	body     []byte // 消息体
	checkSum uint32 // 校验和
}

func (p Protocol) WriteTo(w io.Writer) error {
	// write ver
	err := binary.Write(w, binary.BigEndian, p.ver)
	if err != nil {
		return err
	}
	// write type
	err = binary.Write(w, binary.BigEndian, p.typ)
	if err != nil {
		return err
	}
	// write length
	err = binary.Write(w, binary.BigEndian, p.len)
	if err != nil {
		return err
	}
	// data
	n, err := w.Write(p.body)
	if err != nil {
		return err
	}
	if n != int(p.len) {
		return ErrUnexpectedLength
	}
	// checksum
	p.checkSum = adler32.Checksum(p.body)
	err = binary.Write(w, binary.BigEndian, p.checkSum)
	if err != nil {
		return err
	}
	//log.Debugf("write %d bytes, data:%s", p.len, p.body)
	return err
}

func (p *Protocol) ReadFrom(w io.Reader) error {
	// read ver
	err := binary.Read(w, binary.BigEndian, &p.ver)
	if err != nil {
		return err
	}
	// read type
	err = binary.Read(w, binary.BigEndian, &p.typ)
	if err != nil {
		return err
	}
	// length
	err = binary.Read(w, binary.BigEndian, &p.len)
	if err != nil {
		return err
	}
	// read body
	var body = make([]byte, p.len)
	n, err := io.LimitReader(w, int64(p.len)).Read(body)
	if n != int(p.len) {
		return ErrUnexpectedLength
	}
	p.body = body
	// checksum
	checkSum := adler32.Checksum(p.body)
	err = binary.Read(w, binary.BigEndian, &p.checkSum)
	if err != nil {
		return err
	}
	if checkSum != p.checkSum {
		return ErrUnexpectedCheckSum
	}
	//log.Debugf("read %d bytes, data:%s", p.len, p.body)
	return err
}

func WriteJson(w io.Writer, data interface{}) error {
	msgBody, err := json.Marshal(data)
	if err != nil {
		return err
	}
	msg := Protocol{ver: Ver, typ: JsonMsg, len: int32(len(msgBody)), body: msgBody}
	return msg.WriteTo(w)
}
