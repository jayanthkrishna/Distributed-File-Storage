package p2p

import "net"

// Peer is an interface that represents the remote node
type Peer interface {
	net.Conn
	Send([]byte) error
}

// Transport is anything that handles the communication between nodes in the network.
// Can be of anything(TCP, Web Sockets, UDP)
type Transport interface {
	ListenAndAccept() error
	Consume() <-chan RPC
	Dial(string) error
	Close() error
}
