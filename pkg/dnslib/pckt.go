package dnslib

// DNS packet that is sent across a transport
type DNSPacket struct {
	Header    DNSHdr
	Questions []DNSQuestion
	Answers   []DNSRecord
	// Authorities []DNSRecord
	// Resources   []DNSRecord
}

func (dp *DNSPacket) Unmarshall(msg []byte) error {
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
