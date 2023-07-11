package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	hubrpc "github.com/Mihalic2040/Hub-rpc"
	"github.com/Mihalic2040/Hub-rpc/src/proto/api"
	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/libp2p/go-libp2p/p2p/muxer/mplex"
	"github.com/libp2p/go-libp2p/p2p/muxer/yamux"
	quic "github.com/libp2p/go-libp2p/p2p/transport/quic"
	"github.com/libp2p/go-libp2p/p2p/transport/tcp"
	"github.com/libp2p/go-libp2p/p2p/transport/websocket"
	"github.com/multiformats/go-multiaddr"
)

var endpoint = "/hub/0.0.1"

var port = "0"
var ip = "0.0.0.0"

func Myhandler(input *api.Request) (response api.Response, err error) {

	return api.Response{Payload: input.Payload, Status: 200}, nil
}

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

	handlers := hubrpc.HandlerMap{}

	handlers.HandleFunc("/", Myhandler)

	host.SetStreamHandler(protocol.ID(endpoint), func(stream network.Stream) {
		hubrpc.Stream_handler(stream, handlers)
	})

	kademliaDHT, err := dht.New(context.Background(), host, dht.Mode(dht.ModeAutoServer))
	if err != nil {
		panic(err)
	}

	kademliaDHT.Bootstrap(context.Background())

	log.Println(host.Addrs())

	sourceMultiAddrcnn, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/127.0.0.1/tcp/7778/ws/p2p/12D3KooWAipT8QmimC7WWnhRBcyjaC7geb7Z9ynygkFHJyPUyHKc"))
	if err != nil {
		log.Println("[DHT:Bootstrap] Fail to parse multiaddr: ", err)
	}

	peerinfo, err := peer.AddrInfoFromP2pAddr(sourceMultiAddrcnn)

	host.Connect(context.Background(), *peerinfo)

	log.Println(host.Peerstore().Peers().String())
	for {

	}
}
