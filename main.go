package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/Ghvstcode/DNS/pkg/dnslib"
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

func main() {
	input := "response_packet.txt"

	if len(os.Args) > 1 {
		input = os.Args[1]
	}

	data, err := os.ReadFile(input)
	if err != nil {
		log.Printf("error reading file %v", err)
		return
	}

	dp := &dnslib.DNSPacket{}
	if err := dp.Unmarshall(data); err != nil {
		panic(err)
	}

	if err := json.NewEncoder(os.Stdout).Encode(dp); err != nil {
		panic(err)
	}
}
