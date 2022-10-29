package dnslib

import (
	"fmt"
)

// DNS packet that is sent across a transport
type DNSPacket struct {
	Header    DNSHdr
	Questions []DNSQuestion
	Answers   []DNSRecord
	// Authorities []DNSRecord
	// Resources   []DNSRecord
}

//{ID:41598 QR:1 Opcode:0 AA:0 TC:0 RD:1 RA:1 Z:0 RCODE:0 QDCOUNT:1 ANCOUNT:1 NSCOUNT:0 ARCOUNT:0}
//{id:41598 recursionDesired:true truncatedMessage:false authoritativeAnswer:false opcode:0 response:true rescode:0 checkingDisabled:false authedData:false z:false recursionAvailable:true QdCount:1 AnCount:1 NsCount:0 ArCount:0 QR:false}
func (dp *DNSPacket) Unmarshall(msg []byte) error {
	fmt.Println("LENMSGALL", len(msg))
	var off int
	var err error
	if err := dp.Header.unmarshall(msg[0:12]); err != nil {
		return err
	}

	// Header size is 12 bytes
	off = 12

	dp.Questions = make([]DNSQuestion, int(dp.Header.QdCount))
	for i := 0; i < int(dp.Header.QdCount); i++ {
		newOff, err := dp.Questions[i].unmarshall(msg[off:])
		if err != nil {
			return err
		}
		fmt.Println("AFTER FIRST QUEST OFF", off)
		off += newOff
	}

	dp.Answers = make([]DNSRecord, int(dp.Header.AnCount))
	for idx := range dp.Answers {
		newOff, err := dp.Answers[idx].unmarshall(msg[off:])
		if err != nil {
			return err
		}

		off += newOff
	}

	return err
}
