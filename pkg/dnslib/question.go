package dnslib

import "fmt"

type DNSQuestion struct {
	QNAME  string
	QTYPE  QTYPE
	QCLASS QCLASS
}

type QCLASS uint16

// https://www.iana.org/assignments/dns-parameters/dns-parameters.xhtml#dns-parameters-2
const (
	QClassReserved QCLASS = iota
	QClassIN
	QClassCSNET
	QClassCHAOS
	QClassHESIOD
	QClassNONE
	QClassANY
)

type QTYPE uint16

const (
	QTypeUnkown QTYPE = iota
	// A host address
	QTypeA
	// An authoritative name server
	QTypeNS
	// A mail destination
	QTypeMD
	// A mail forwarder
	QTypeMF
	// The canonical name for an alias
	QTypeCNAME
	// Marks the start of a zone of authority
	QTypeSOA
	// A mailbox domain name
	QTypeMB
	// A mail group member
	QTypeMG
	// A mail rename domain name
	QTypeMR
	// A null RR
	QTypeNULL
	// A well known service description
	QTypeWKS
	// A domain name pointer
	QTypePTR
	// Host information
	QTypeHINFO
	// Mailbox or mail list information
	QTypeMINFO

	// A request for a transfer of an entire zone of authority
	QTypeAXFR QTYPE = 252
	// A request for mailbox-related records (MB, MG or MR)
	QTypeMAILB QTYPE = 253
	// A request for mail agent RRs (MD and MF)
	QTypeMAILA QTYPE = 254
	QTypeAll   QTYPE = 255
)

func (dq *DNSQuestion) unmarshall(msg []byte) (off int, err error) {
	// var off int
	// var err error
	dq.QNAME, off, err = readQname(msg, off)
	if err != nil {
		return off, err
	}

	// Read the Question type
	qt, err := ReadUint16(msg, off)
	dq.QTYPE = QTYPE(qt)
	if err != nil {
		return off, err
	}
	fmt.Println("-78QTYPE", dq.QTYPE, qt)
	off += 2

	// Read the Question Class
	qc, err := ReadUint16(msg, off)
	dq.QCLASS = QCLASS(qc)
	if err != nil {
		return off, err
	}
	off += 2

	return off, err
}
