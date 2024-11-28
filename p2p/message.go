package p2p

import "net"

// Message Holds the data that is being sent between the nodes in the network
type RPC struct {
	From    net.Addr
	Payload []byte
}
