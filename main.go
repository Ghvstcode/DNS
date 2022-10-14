package main

import (
	"errors"
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

func numToRcode(num int) RCODE {
	return RCODE(num)
}

// DNSQuestion represents the question section which is used to carry the "question" in most queries,
// i.e., the parameters that define what is being asked
type DNSQuestion struct {
	// RName is a domain name represented as a sequence of labels, where each label consists of a length octet
	// followed by that number of octets.The domain name terminates with the
	// zero length octet for the null label of the root
	QNAME string
	// QTYPE specifies the type of the query
	QTYPE uint16
	// QCLASS  a two octet code that specifies the class of the query
	QCLASS uint16
}

type DNSRecord struct {
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

// DNS packet that is sent across a transport
type DNSPacket struct {
	// Header of the DNS packet that contains info about the packet. It is always present
	Header DNSHdr
	// Questions The question section is used to carry the "question" in most queries,
	// i.e., the parameters that define what is being asked
	Questions []DNSQuestion
	// The Answer, Authorities, and additional sections all share the same
	// format: a variable number of resource records, where the number of
	// records is specified in the corresponding count field in the header
	Answers     []DNSRecord
	Authorities []DNSRecord
	Resources   []DNSRecord
}

func (dq *DNSQuestion) unmarshall(msg []byte) error {}
