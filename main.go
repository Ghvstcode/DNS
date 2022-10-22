package main

import (
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
	fmt.Println("Contents of file:")
	fmt.Println(string(data))

	dp := dnslib.DNSPacket{}
	fmt.Println("DNS PACKET1", dp)
	if err := dp.Unmarshall(data); err != nil {
		panic(err)
	}
	fmt.Println("DNS PACKET", dp)
}
