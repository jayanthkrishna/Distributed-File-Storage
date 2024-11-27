package p2p

// Peer is an interface that represents the remote node
type Peer interface {
}

// Transport is anything that handles the communication between nodes in the network.
// Can be of anything(TCP, Web Sockets, UDP)
type Transport interface {
	ListenAndAccept() error
}
