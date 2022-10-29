package main

import (
	"encoding/json"
	"fmt"
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
	data, err := os.ReadFile("response_packet.txt")
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}

	fmt.Println("INCLEN", len(data))
	dp := &dnslib.DNSPacket{}
	if err := dp.Unmarshall(data); err != nil {
		panic(err)
	}

	if err := json.NewEncoder(os.Stdout).Encode(dp); err != nil {
		panic(err)
	}
}
