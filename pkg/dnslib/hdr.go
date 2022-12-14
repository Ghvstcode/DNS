package dnslib

import (
	"errors"
)

type RCODE int

// DNS MSG HEADER
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

type DNSHdr struct {
	id                  uint16 // 16 bits
	recursionDesired    bool   // 1 bit
	truncatedMessage    bool   // 1 bit
	authoritativeAnswer bool   // 1 bit
	opcode              uint8  // 4 bits
	response            bool   // 1 bit
	rescode             RCODE  // 4 bits
	checkingDisabled    bool   // 1 bit
	authedData          bool   // 1 bit
	z                   bool   // 1 bit
	recursionAvailable  bool   // 1 bit
	QdCount             uint16 // 16 bits
	AnCount             uint16 // 16 bits
	NsCount             uint16 // 16 bits
	ArCount             uint16 // 16 bits
	QR                  bool
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
	// Read Header Flags
	dh.unmarshallHeaderFlags(bits)
	off += 2
	dh.QdCount, err = ReadUint16(msg, off)
	if err != nil {
		return err
	}
	off += 2
	dh.AnCount, err = ReadUint16(msg, off)
	if err != nil {
		return err
	}
	off += 2
	dh.NsCount, err = ReadUint16(msg, off)
	if err != nil {
		return err
	}
	off += 2
	dh.ArCount, err = ReadUint16(msg, off)
	if err != nil {
		return err
	}
	return nil
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

func (dh *DNSHdr) marshallHeader(msg []byte) ([]byte, error) {
	off := 1
	err := PutUint16(msg, off, dh.id)
	if err != nil {
		return nil, err
	}
	// Increment the offset after the previous uint16 read
	off += 2
	err = PutUint16(msg, off, dh.getBits())
	if err != nil {
		return nil, err
	}

	off += 2
	err = PutUint16(msg, off, dh.QdCount)
	if err != nil {
		return nil, err
	}
	off += 2
	err = PutUint16(msg, off, dh.AnCount)
	if err != nil {
		return nil, err
	}
	off += 2
	err = PutUint16(msg, off, dh.NsCount)
	if err != nil {
		return nil, err
	}
	off += 2
	err = PutUint16(msg, off, dh.ArCount)
	if err != nil {
		return nil, err
	}
	// TODO: pack flag as well.
	return msg, nil
}

func (dh *DNSHdr) getBits() uint16 {
	var bits uint16
	bits = uint16(dh.opcode)<<11 | uint16(dh.rescode&0xF)
	if dh.response {
		bits |= 1 << 0xF
	}
	if dh.authoritativeAnswer {
		bits |= 1 << 0xA
	}
	if dh.truncatedMessage {
		bits |= 1 << 0x9
	}
	if dh.recursionDesired {
		bits |= 1 << 0x8
	}

	if dh.recursionAvailable {
		bits |= 1 << 0x7
	}

	if dh.z {
		bits |= 1 << 0x6
	}
	if dh.authedData {
		bits |= 1 << 0x5
	}
	if dh.checkingDisabled {
		bits |= 1 << 0x4
	}

	return bits
}

func ReadUint16(msg []byte, off int) (uint16, error) {
	if off+2 > len(msg) {
		return 0, newErr("overflow of slice while reading uint from message")
	}
	// https://cs.opensource.google/go/go/+/refs/tags/go1.19.1:src/encoding/binary/binary.go;l=140
	return uint16(msg[off+1]) | uint16(msg[off])<<0x8, nil
}

func PutUint16(msg []byte, off int, v uint16) error {
	if off+2 > len(msg) {
		return newErr("overflow of slice while packing uint16")
	}
	// https: // cs.opensource.google/go/go/+/refs/tags/go1.19.1:src/encoding/binary/binary.go;l=143
	msg[off-1] = byte(v >> 8)
	msg[off] = byte(v)

	return nil
}

// This function is a stripped down version of this https://github.com/miekg/dns/blob/master/msg.go#L378
func readQname(msg []byte, off int) (string, int, error) {
	off1 := 0
	lenmsg := len(msg)
	ptr := 0 // number of pointers followed
	res := make([]byte, 0)
Loop:
	for {
		if off >= lenmsg {
			return "", lenmsg, newErr("overflow of slice")
		}
		c := int(msg[off])
		off++
		switch c & 0xC0 {
		case 0x00:
			if c == 0x00 {
				break Loop
			}

			// literal string
			if off+c > lenmsg {
				return "", lenmsg, newErr("overflow of slice")
			}

			for _, b := range msg[off : off+c] {
				if isDomainNameLabelSpecial(b) {
					res = append(res, '\\', b)
				} else {
					res = append(res, b)
				}
			}
			res = append(res, '.')
			off += c
		case 0xC0:
			if off >= lenmsg {
				return "", lenmsg, newErr("overflow of slice")
			}
			c1 := msg[off]
			off++
			if ptr == 0 {
				off1 = off
			}

			ptr++
			off = (c^0xC0)<<8 | int(c1)
		default:
			// 0x80 and 0x40 are reserved
			return "", lenmsg, newErr("bad domain name")
		}
	}
	if ptr == 0 {
		off1 = off
	}
	if len(res) == 0 {
		return ".", off1, nil
	}
	return string(res), off1, nil
}

func isDomainNameLabelSpecial(b byte) bool {
	switch b {
	case '.', ' ', '\'', '@', ';', '(', ')', '"', '\\':
		return true
	}
	return false
}
