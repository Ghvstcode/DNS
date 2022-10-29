package dnslib

import (
	"encoding/binary"
)

type DNSRecord struct {
	QNAME    string
	QTYPE    QTYPE
	QCLASS   QCLASS
	TTL      uint32
	RDLENGTH uint16
	RDATA    []byte
}

func (dr *DNSRecord) unmarshall(msg []byte) (int, error) {
	var off int
	var err error

	//dr.QNAME, off, err = readQname(msg, off)
	//if err != nil {
	//	panic(err)
	//	return 0, err
	//}

	off += 2

	// Read the Question type
	qt, err := ReadUint16(msg, off)
	dr.QTYPE = QTYPE(qt)
	if err != nil {
		return 0, err
	}
	off += 2

	// Read the Question Class
	qc, err := ReadUint16(msg, off)
	dr.QCLASS = QCLASS(qc)
	if err != nil {
		return 0, err
	}
	off += 2

	dr.TTL = binary.BigEndian.Uint32(msg[off:])
	off += 4

	// Read RDLength
	dr.RDLENGTH, err = ReadUint16(msg, off)
	if err != nil {
		return 0, err
	}
	off += 3

	finalOff := uint16(off) + dr.RDLENGTH

	dr.RDATA = msg[off:finalOff]
	return int(finalOff), nil
}
