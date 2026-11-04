package p2p

import "errors"

//ErrInvalidHandshake is returned if the handshake between
//local and the remote is not established

var ErrInvalidHandshake = errors.New("Invalid Handshake")

type HandshakeFunc func(Peer) error

func NOPHandshakeFunc(Peer) error { return nil }

