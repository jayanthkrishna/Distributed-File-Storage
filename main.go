package main

import (
	"fmt"
	"log"

	"github.com/jayanthkrishna/Distributed-File-Storage/p2p"
)

func main() {
	fmt.Println("hello World")

	tr := p2p.NewTCPTransport(":3000")

	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	// select {}

}
