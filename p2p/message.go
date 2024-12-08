package p2p

// Message Holds the data that is being sent between the nodes in the network
type RPC struct {
	From    string
	Payload []byte
}
