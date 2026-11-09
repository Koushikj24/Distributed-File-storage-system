package p2p

import "net"

// Peer is an interface that represents a remote node
type Peer interface {
	Send([]byte) error
	RemoteAddr() net.Addr
	Close() error
}

// Transport handles communication between thr nodes in the network
// This can be of TCP or UDP or websockets...etc
type Transport interface {
	Dial(string) error
	ListenAndAccept() error 
	Consume() <- chan RPC
	Close() error
}