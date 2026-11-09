package p2p

import (
	"fmt"
	"log"
	"net"
	"errors"
)

type TCPPeer struct {

	// conn is underlying connection of the peer
	conn net.Conn

	// if we dial and retrive a conn =>outbound = true
	// if we accept and retrive a conn =>outbound = false
	outbound bool   
}


func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer{
  return &TCPPeer{
	conn : conn,
	outbound: outbound,
  }
}

func (p *TCPPeer) Send(b []byte) error {
	_, err := p.conn.Write(b)
	return err
}

// RemoteAddr implements the Peer interface and will return the remote address of its underlying connection.
func(p *TCPPeer) RemoteAddr() net.Addr {
	return p.conn.RemoteAddr()
}

func (p *TCPPeer) Close() error{
	return p.conn.Close()
}

type TCPTransportOpts struct {
	 ListenAddr string
	 HandshakeFunc HandshakeFunc
	 Decoder Decoder
	 OnPeer func(Peer) error
}

type TCPTransport struct {
    TCPTransportOpts
	listener net.Listener
    rpcch chan RPC


}



func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
     return &TCPTransport{
		TCPTransportOpts: opts,
		rpcch : make(chan RPC),
	 }
}

// Consume implements the Transport interface, which will return read-only channel
// for reading the incoming messages received from another peer in the network.
func(t *TCPTransport) Consume() <-chan RPC{
	return t.rpcch
}

// Close implements the transport interface.
func(t *TCPTransport) Close() error {
	return t.listener.Close()
}

// Dial implements the transport interface
func(t *TCPTransport) Dial(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	go t.handleConn(conn, true)

	return nil
}

func(t *TCPTransport) ListenAndAccept() error{

	var err error 

	t.listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil{
		 return err
	}
    
	go t.startAcceptLoop()

	log.Printf("TCP Transport listening on port : %s\n", t.ListenAddr)

    return nil
}


func ( t *TCPTransport) startAcceptLoop(){
	for {
		conn, err := t.listener.Accept()

		if errors.Is(err, net.ErrClosed){
			return
		}

		if err != nil{
			fmt.Printf("TCP Accept error : %s\n", err)
		}
	    fmt.Printf("New incoming connection: %s\n", conn.RemoteAddr())

        go t.handleConn(conn, false);

	}
}



func (t *TCPTransport) handleConn(conn net.Conn, outbound bool) {

	var err error

	defer func(){
		fmt.Printf("dropping peer connection: %s", err)
		conn.Close()
	}()

    peer := NewTCPPeer(conn,outbound)
    
    if err := t.HandshakeFunc(peer); err != nil{
		return
	}
    
	if t.OnPeer != nil {
		 if err = t.OnPeer(peer); err!= nil {
			return 
		 }
	}

    // Read loop
	rpc := RPC{}
	for {
	  if err := t.Decoder.Decode(conn, &rpc); err!=nil {
		fmt.Printf("TCP error: %s\n", err)
		continue
	  }
	  rpc.From = conn.RemoteAddr()
	  t.rpcch <- rpc 
	} 
}

