package p2p

import "net"

//RPC holds any represents any arbitrary data
//that is sent over the each transport between the nodes i the network.
type RPC struct {
	 From net.Addr
     Payload []byte
}