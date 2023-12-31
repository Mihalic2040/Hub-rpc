package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/p2p/muxer/mplex"
	"github.com/libp2p/go-libp2p/p2p/muxer/yamux"
	quic "github.com/libp2p/go-libp2p/p2p/transport/quic"
	"github.com/libp2p/go-libp2p/p2p/transport/tcp"
	"github.com/libp2p/go-libp2p/p2p/transport/websocket"
	"github.com/multiformats/go-multiaddr"
)

var endpoint = "/hub/0.0.1"

var port = "7777"
var ip = "0.0.0.0"

func main() {
	log.Println("Starting server...")

	sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%s/", ip, port))
	sourceMultiAddrQuic, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/udp/%s/quic", ip, port))
	var sourceMultiAddrWs multiaddr.Multiaddr
	if port != "0" {
		ports, _ := strconv.Atoi(port)
		portNumber := strconv.Itoa(ports + 1)
		sourceMultiAddrWs, _ = multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%s/ws", ip, portNumber))
	} else {
		sourceMultiAddrWs, _ = multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%s/ws", ip, port))
	}

	addrs := libp2p.ListenAddrStrings(
		sourceMultiAddr.String(),
		sourceMultiAddrQuic.String(),
		sourceMultiAddrWs.String(),
	)

	taranspors := libp2p.ChainOptions(
		libp2p.Transport(tcp.NewTCPTransport),
		libp2p.Transport(quic.NewTransport),

		//For js nodes
		libp2p.Transport(websocket.New),
	)

	muxers := libp2p.ChainOptions(
		libp2p.Muxer("/mplex/", mplex.DefaultTransport),

		//For js nodes
		libp2p.Muxer("/yamux/", yamux.DefaultTransport),
	)

	host, err := libp2p.New(
		taranspors,
		muxers,
		addrs,

		libp2p.EnableHolePunching(),
		libp2p.NATPortMap(),
		libp2p.EnableNATService(),
		libp2p.EnableRelayService(),
	)
	if err != nil {
		panic(err)
	}

	//register rpc server

	kademliaDHT, err := dht.New(context.Background(), host, dht.Mode(dht.ModeAutoServer))
	if err != nil {
		panic(err)
	}

	kademliaDHT.Bootstrap(context.Background())

	log.Println(host.Addrs())
	log.Println(host.ID().String())

	for {
	}
}
