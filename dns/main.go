package main

import (
	"fmt"
	dnsbuf "go-dns/dns_buf"
	"io"
	"log"
	"os"
)

func main() {

	f, err := os.Open("response_packet.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	buf := dnsbuf.New()
	buf.Load(data)

	packet, err := dnsbuf.ReadPacket(buf)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", packet.Header)

	for _, q := range packet.Questions {
		fmt.Printf("%+v\n", q)
	}
	for _, rec := range packet.Answers {
		fmt.Println(rec)
	}
	for _, rec := range packet.Authorities {
		fmt.Println(rec)
	}
	for _, rec := range packet.Resources {
		fmt.Println(rec)
	}
}
