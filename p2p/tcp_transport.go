package p2p

import (
	"fmt"
	"log"
	"net"
)

// TCP Peer represents the remote node over a TCP connection
type TCPPeer struct {
	conn net.Conn
	// If we dial a connection , it is outbound. If we accept a conn, it is inbound
	outbound bool
}

type TCPTransportOpts struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
	OnPeer        func(Peer) error
}

type TCPTransport struct {
	TCPTransportOpts

	listener net.Listener
	rpcch    chan RPC
}

func NewTCPPeer(conn net.Conn, outbound bool) Peer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

func (p *TCPPeer) Close() error {
	return p.conn.Close()
}
func NewTCPTransport(opts TCPTransportOpts) Transport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		rpcch:            make(chan RPC),
		// listener:      &net.TCPListener{},
		// mu:            sync.RWMutex{},
		// peers:         make(map[net.Addr]Peer),
	}
}

// Consume implements the transport interface which will return a r4ead only channel.
// for reading the incoming messages received from another peer in the network
func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcch
}
func (t *TCPTransport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.ListenAddr)

	if err != nil {
		return err
	}

	go t.startLoop()

	log.Printf("TCP transport listening on port: %s\n", t.ListenAddr)

	return nil

}

func (t *TCPTransport) startLoop() {
	for {
		conn, err := t.listener.Accept()

		if err != nil {
			fmt.Printf("TCP accept error: %s\n", err)
			continue
		}

		go t.handleConn(conn)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn) {

	var err error
	defer func() {
		fmt.Printf("dropping peer connection : %S\n", err)
		conn.Close()
	}()

	peer := NewTCPPeer(conn, true)
	if err = t.HandshakeFunc(peer); err != nil {
		conn.Close()
		return

	}

	if t.OnPeer != nil {
		if err = t.OnPeer(peer); err != nil {
			return
		}
	}

	//  Read Loop
	msg := &RPC{}

	for {
		if err := t.Decoder.Decode(conn, msg); err != nil {

			fmt.Printf("TCP Error : %s\n", err)
			continue

		}
		msg.From = conn.RemoteAddr()
		t.rpcch <- *msg

	}

}
