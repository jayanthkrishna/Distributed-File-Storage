package p2p

import (
	"fmt"
	"net"
	"sync"
)

// TCP Peer represents the remote node over a TCP connection
type TCPPeer struct {
	conn net.Conn
	// If we dial a connection , it is outbound. If we accept a conn, it is inbound
	outbound bool
}

type TCPTransport struct {
	listenAddress string
	listener      net.Listener
	handshake     HandshakeFunc
	mu            sync.RWMutex
	peers         map[net.Addr]Peer
}

func NewTCPPeer(conn net.Conn, outbound bool) Peer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

func NewTCPTransport(listenAddr string) Transport {
	return &TCPTransport{
		listenAddress: listenAddr,
		// listener:      &net.TCPListener{},
		// mu:            sync.RWMutex{},
		// peers:         make(map[net.Addr]Peer),
	}
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.listenAddress)

	if err != nil {
		return err
	}

	t.startLoop()

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

type Temp struct {
}

func (t *TCPTransport) handleConn(conn net.Conn) {

	// peer := NewTCPPeer(conn, true)

	if err := t.handshake(conn); err != nil {

	}

	//  Read Loop

	// defer conn.Close()
	fmt.Printf("New incoming connection: %+v\n", conn)
}
