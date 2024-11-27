package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	listenAddr := ":4000"
	tr := NewTCPTransport(listenAddr)

	trr := tr.(*TCPTransport)

	assert.Equal(t, trr.listenAddress, listenAddr)

	assert.Nil(t, tr.ListenAndAccept())
}
