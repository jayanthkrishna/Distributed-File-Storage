package main

import (
	"bytes"
	"fmt"

	"github.com/jayanthkrishna/Distributed-File-Storage/p2p"
)

func makeServer(listenAddr string, nodes ...string) *FileServer {
	tcpOpts := p2p.TCPTransportOpts{
		ListenAddr:    listenAddr,
		HandshakeFunc: p2p.NOPHandShakeFunc,
		Decoder:       p2p.DefaultDecoder{},
	}

	tr := p2p.NewTCPTransport(tcpOpts)

	fileServeropts := FileServerOpts{
		StorageRoot:       listenAddr + "_network",
		PathTransformFunc: CASPathTransformFunc,
		Transport:         tr,
		BootstrapNodes:    nodes,
	}

	fmt.Println(listenAddr)
	s := NewFileServer(fileServeropts)

	tr.OnPeer = s.OnPeer

	return s

}

func main() {
	fmt.Println("hello World")

	s1 := makeServer(":3000", "")

	s2 := makeServer(":4000", ":3000")

	go s1.Start()
	go s2.Start()

	data := bytes.NewReader([]byte("My big data file is here!!!"))

	s2.Store("myprivatekey", data)
	select {}
}
