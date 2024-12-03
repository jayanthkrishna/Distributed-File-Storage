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

	fileServeropts := FileServerOpts{
		StorageRoot:       "3000_network",
		PathTransformFunc: CASPathTransformFunc,
		Transport:         tr,
	}

	s := NewFileServer(fileServeropts)

	if err := s.Start(); err != nil {
		log.Fatal(err)

	}

	select {}
}
