package p2p

// Peer is an interface that represents a remote node
type Peer interface {
	Close() error
}

// Transport handles communication between thr nodes in the network
// This can be of TCP or UDP or websockets...etc
type Transport interface {
	ListenAndAccept() error 
	Consume() <- chan RPC
}