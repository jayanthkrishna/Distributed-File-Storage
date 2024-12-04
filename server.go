package main

import (
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
	peer     map[string]p2p.Peer
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
		peer:           make(map[string]p2p.Peer),
	}
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
		case msg := <-s.Transport.Consume():
			fmt.Println(msg)
		case <-s.quitch:
			return

		}
	}
}

func (s *FileServer) OnPeer(p p2p.Peer) error {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()

	s.peer[p.RemoteAddr().String()] = p

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
	return s.store.Write(key, r)
}
