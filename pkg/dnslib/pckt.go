package dnslib

import "fmt"

// DNS packet that is sent across a transport
type DNSPacket struct {
	Header      DNSHdr
	Questions   []DNSQuestion
	Answers     []DNSRecord
	Authorities []DNSRecord
	Resources   []DNSRecord
}

func (dp *DNSPacket) Unmarshall(msg []byte) error {
	var off int
	var err error
	if err := dp.Header.unmarshall(msg[0:12]); err != nil {
		return err
	}
	// Header size is 12 bytes
	off = 12

	dp.Questions = make([]DNSQuestion, int(dp.Header.QdCount)-1)
	for idx := range dp.Questions {
		off, err = dp.Questions[idx].unmarshall(msg[off:])
		if err != nil {
			return err
		}
	}

	fmt.Println("int(dp.Header.AnCount)", int(dp.Header.AnCount))
	dp.Answers = make([]DNSRecord, int(dp.Header.AnCount))
	for idx := range dp.Answers {
		fmt.Println("HIT-34")
		off, err = dp.Answers[idx].unmarshall(msg[off:])
		if err != nil {
			fmt.Println("ERR", err)
			return err
		}
		fmt.Println("dp.Answers[idx]", dp.Answers[idx].QTYPE)
	}

	for idx := range dp.Authorities {
		off, err = dp.Answers[idx].unmarshall(msg[off:])
		if err != nil {
			return err
		}
	}

	for idx := range dp.Resources {
		off, err = dp.Answers[idx].unmarshall(msg[off:])
		if err != nil {
			return err
		}
	}

	return err
}
