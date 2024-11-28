package main

import (
	"fmt"
	"log"

	"github.com/jayanthkrishna/Distributed-File-Storage/p2p"
)

func main() {
	fmt.Println("hello World")

	tcpOpts := p2p.TCPTransportOpts{
		ListenAddr:    ":3000",
		HandshakeFunc: p2p.NOPHandShakeFunc,
		Decoder:       p2p.DefaultDecoder{},
	}

	tr := p2p.NewTCPTransport(tcpOpts)

	fmt.Println("here - 1")
	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("here - 2")
	go func() {
		fmt.Println("here - in routine")
		for {
			msg := <-tr.Consume()

			fmt.Printf("%s\n", string(msg.Payload))
		}

	}()
	fmt.Println("Server started on port 3000")
	select {}

}
