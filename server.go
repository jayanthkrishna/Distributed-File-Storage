package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/jayanthkrishna/Distributed-File-Storage/p2p"
)

type FileServerOpts struct {
	StorageRoot       string
	PathTransformFunc PathTransformFunc
	Transport         p2p.Transport
	BootstrapNodes    []string
}

type FileServer struct {
	FileServerOpts
	peerLock sync.Mutex
	peers    map[string]p2p.Peer
	store    *Store
	quitch   chan struct{}
}

func NewFileServer(opts FileServerOpts) *FileServer {
	storeOpts := StoreOpts{
		Root:              opts.StorageRoot,
		PathTransformFunc: opts.PathTransformFunc,
	}
	return &FileServer{
		FileServerOpts: opts,
		store:          NewStore(storeOpts),
		quitch:         make(chan struct{}),
		peers:          make(map[string]p2p.Peer),
	}
}

type Message struct {
	From    string
	Payload any
}
type Payload struct {
	key  string
	Data []byte
}

func (s *FileServer) handleMessage(p *Payload) error {
	return nil

}
func (s *FileServer) broadcast(p *Payload) error {
	peers := []io.Writer{}
	for _, peer := range s.peers {
		peers = append(peers, peer)

	}

	mw := io.MultiWriter(peers...)

	return gob.NewEncoder(mw).Encode(p)
}

func (s *FileServer) bootstrapNetwork() error {
	for _, addr := range s.BootstrapNodes {
		if len(addr) == 0 {
			continue
		}
		go func(addr string) {
			if err := s.Transport.Dial(addr); err != nil {
				log.Println("Dial Error : ", err)
				return
			}
			fmt.Println("Bootstraped Addr ", addr)
		}(addr)

	}
	return nil
}
func (s *FileServer) loop() {

	defer func() {
		log.Println("File server stopped")
		s.Transport.Close()
	}()
	for {
		select {
		case rpc := <-s.Transport.Consume():
			var msg Message
			if err := gob.NewDecoder(bytes.NewReader(rpc.Payload)).Decode(&msg); err != nil {
				log.Println("decoding error: ", err)
			}

			peer, ok := s.peers[rpc.From]

			if !ok {
				panic("Peer not found in the map.")
			}
			b := make([]byte, 1000)
			if _, err := peer.Read(b); err != nil {
				panic("read error")
			}
			fmt.Printf("received: %s\n", string(b))

			peer.(*p2p.TCPPeer).Wg.Done()
		case <-s.quitch:
			return

		}
	}
}

func (s *FileServer) OnPeer(p p2p.Peer) error {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()

	s.peers[p.RemoteAddr().String()] = p

	log.Printf("connected with remote %s", p.RemoteAddr())

	return nil

}
func (s *FileServer) Stop() {
	close(s.quitch)
}
func (s *FileServer) Start() error {

	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}
	fmt.Println("After listen and accept")
	s.bootstrapNetwork()
	s.loop()
	return nil
}

func (s *FileServer) Store(key string, r io.Reader) error {

	buf := new(bytes.Buffer)

	msg := Message{
		Payload: []byte("storagekey"),
	}

	if err := gob.NewEncoder(buf).Encode(msg); err != nil {
		return err
	}

	for _, peer := range s.peers {
		if err := peer.Send(buf.Bytes()); err != nil {
			return err
		}
	}

	payload := []byte("This is a big file!!!")
	for _, peer := range s.peers {
		if err := peer.Send(payload); err != nil {
			return err
		}
	}
	return nil

	// tee := io.TeeReader(r, buf)

	// if err := s.store.Write(key, tee); err != nil {
	// 	return err
	// }

	// p := &Payload{
	// 	key:  key,
	// 	Data: buf.Bytes(),
	// }

	// fmt.Println(buf.Bytes())
	// return s.broadcast(p)
}
