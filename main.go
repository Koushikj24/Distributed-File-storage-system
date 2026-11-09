package main

import (
	"log"

	"github.com/Koushikj24/foreverstore/p2p"
)

// import (
//     "log"
// 	"fmt"
// 	"github.com/Koushikj24/foreverstore/p2p"
// )
// func OnPeer(p2p.Peer) error {
// 	fmt.Println("Working on a logic with peer outside TCPTransport")
// 	return nil
// }

// func main(){
// 	tcpOpts := p2p.TCPTransportOpts{
// 		ListenAddr: ":3000",
// 		HandshakeFunc: p2p.NOPHandshakeFunc,
// 		Decoder: p2p.DefaultDecoder{},
// 		OnPeer: OnPeer,
// 	}

// 	tr:= p2p.NewTCPTransport(tcpOpts)

// 	go func(){
// 		for {
// 			msg:=  <-tr.Consume()
// 			fmt.Printf("%+v\n", msg)
// 		}
// 	}()
// 	if err := tr.ListenAndAccept(); err != nil {
// 		log.Fatal(err)
// 	}

// 	select{}

// }

func makeServer(listenAddr string, nodes ...string) *FileServer{
	tcptransportOpts := p2p.TCPTransportOpts{
		ListenAddr: listenAddr,
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder: p2p.DefaultDecoder{},
	
	}
	
	tcpTransport := p2p.NewTCPTransport(tcptransportOpts)

	FileServerOpts := FileServerOpts{
		StorageRoot: listenAddr + "_network",
		PathTransformFunc: CASPathTransformFunc,
		Transport: tcpTransport,
		BootstrapNodes: nodes,
	}
	
	s := NewFileServer(FileServerOpts)
	tcpTransport.OnPeer = s.OnPeer

	return s
}

func main(){

	s1 := makeServer(":3000", "")
	s2 := makeServer(":4000", ":3000")

	go func(){
		log.Fatal(s1.Start())
	}()

	s2.Start()

}

