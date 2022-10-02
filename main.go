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

func numToRcode(num int) RCODE {
	return RCODE(num)
}

type DNSHdr struct {
	id                   uint16 // 16 bits
	recursionDesired     bool   // 1 bit
	truncatedMessage     bool   // 1 bit
	authoritativeAnswer  bool   // 1 bit
	opcode               uint8  // 4 bits
	response             bool   // 1 bit
	rescode              RCODE  // 4 bits
	checkingDisabled     bool   // 1 bit
	authedData           bool   // 1 bit
	z                    bool   // 1 bit
	recursionAvailable   bool   // 1 bit
	questions            uint16 // 16 bits
	answers              uint16 // 16 bits
	authoritativeEntries uint16 // 16 bits
	resourceEntries      uint16 // 16 bits
}

type Question struct {
	// RName is a domain name represented as a sequence of labels, where each label consists of a length octet
	// followed by that number of octets.The domain name terminates with the
	// zero length octet for the null label of the root
	QNAME string
	// QTYPE specifies the type of the query
	QTYPE uint16
	// QCLASS  a two octet code that specifies the class of the query
	QCLASS uint16
}

type Record struct {
	// RName a domain name to which this resource record pertains
	RName string
	// RType two octets containing one of the RR type codes
	RType uint16
	// RClass two octets which specify the class of the data
	RClass uint16
	// RTTL is a a 32 bit unsigned integer that specifies the time
	// interval (in seconds) that the resource record may be
	// cached before it should be discarded
	RTTL uint32
	// Rdlength an unsigned 16 bit integer that specifies the length
	Rdlength uint16
}
