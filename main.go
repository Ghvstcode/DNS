package main

import (
	"fmt"
)

//DNS packet structure
//+---------------------+
//|        Header       |
//+---------------------+
//|       Question      | Question for the name server
//+---------------------+
//|        Answer       | records that were requested from above
//+---------------------+
//|      Authority      | Name servers that resolve queries recursively
//+---------------------+
//|      Additional     | additional records

// DNS HEADER
//+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+---------+
//|  ID(16bit) - assigned  randomly                      |
//+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+---------+
//|QR(1bi)|Opcode(4bi)|AA(4bi)|TC|RD|RA|Z    |   RCODE      |
//+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+----------
//|                    QDCOUNT                           |
//+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+---------+
//|                    ANCOUNT                           |
//+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+---------+
//|                    NSCOUNT                           |
//+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+------+
//|                    ARCOUNT                           |
//+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+------+

type BytePacket struct {
	buf []byte
	pos uint16
}

func main() {
}

func NewBytePacket() BytePacket {
	return BytePacket{
		buf: make([]byte, 512),
		pos: 0,
	}
}

func (b *BytePacket) currentPosition() uint16 {
	return b.pos
}

func (b *BytePacket) seek(pos uint16) {
	if pos > 512 {
		b.pos = pos
	}
}

func (b *BytePacket) step(pos uint16) {
	b.pos += pos
}

func (b *BytePacket) read() (byte, error) {
	if b.pos >= 512 {
		return 0, fmt.Errorf("end of buffer")
	}
	b.pos += 1
	return b.buf[b.pos], nil
}

func (b *BytePacket) get(pos uint16) (byte, error) {
	if pos >= 512 {
		return 0, fmt.Errorf("end of buffer")
	}
	return b.buf[pos], nil
}

func (b *BytePacket) getRange(start, len uint16) ([]byte, error) {
	if start+len >= 512 {
		return nil, fmt.Errorf("end of buffer")
	}
	return b.buf[start : start+len], nil
}

// Error represents a DNS error.
type Error struct{ err string }

func (e *Error) Error() string {
	if e == nil {
		return "dns: <nil>"
	}
	return "dns: " + e.err
}

type RCODE int

const (
	NOERROR RCODE = iota
	FORMERR
	SERVFAIL
	NXDOMAIN
	NOTIMP
	REFUSED
)
