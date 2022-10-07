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

func newErr(e string) error {
	return Error{err: e}
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

func (dh *DNSHdr) unmarshall(msg []byte) error {
	// Check to see if the lenght of the message is the correct size! The message size is usually 12 bytes
	if len(msg) != 12 {
		return errors.New("invalid header size")
	}
	var err error

	off := 0

	dh.id, err = ReadUint16(msg, off)
	if err != nil {
		return err
	}
	// Increment the offset after the previous uint16 read
	off += 2
	bits, err := ReadUint16(msg, off)
	if err != nil {
		return err
	}
	off += 2
}

func (dh *DNSHdr) unmarshallHeaderFlags(bits uint16) {
	dh.recursionDesired = (bits & (1 << 0x8)) != 0
	dh.response = (bits & (1 << 0xF)) != 0
	dh.opcode = uint8((bits >> 0xB) & 0xF)
	dh.authoritativeAnswer = (bits & (1 << 0xA)) != 0
	dh.truncatedMessage = (bits & (1 << 0x9)) != 0
	dh.recursionAvailable = (bits & (1 << 0x7)) != 0
	dh.z = (bits & (1 << 0x6)) != 0
	dh.authedData = (bits & (1 << 0x5)) != 0
	dh.checkingDisabled = (bits & (1 << 0x4)) != 0
	dh.opcode = uint8((bits >> 0xB) & 0xF)
}

func ReadUint16(msg []byte, off int) (uint16, error) {
	if off+2 > len(msg) {
		return 0, newErr("overflow of slice while reading uint from message")
	}
	// https://cs.opensource.google/go/go/+/refs/tags/go1.19.1:src/encoding/binary/binary.go;l=140
	return uint16(msg[1]) | uint16(msg[0])<<0x8, nil
}
