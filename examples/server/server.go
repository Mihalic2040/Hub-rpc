package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

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

	return api.Response{}, nil
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

	kademliaDHT, err := dht.New(context.Background(), host)
	if err != nil {
		panic(err)
	}

	kademliaDHT.Bootstrap(context.Background())

	log.Println(host.Peerstore().Peers().String())

	host2, err := libp2p.New(
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

	kademliaDHT2, err := dht.New(context.Background(), host2)
	if err != nil {
		panic(err)
	}

	kademliaDHT2.Bootstrap(context.Background())

	host2Info := peer.AddrInfo{
		ID:    host.ID(),
		Addrs: host.Addrs(),
	}

	log.Println(host2Info)

	if err := host2.Connect(context.Background(), host2Info); err != nil {
		log.Println("Failed to connect host1 to host2:", err)
	}

	log.Println(host2.Peerstore().Peers().String())

	data := api.Request{
		Payload: "gg",
		Handler: "/",
	}

	for {
		time.Sleep(time.Second * 1)
		log.Println(host.ID().String())
		res, err := hubrpc.NewRequest(host.ID().String(), &data, endpoint, *kademliaDHT2, host2)
		if err != nil {
			log.Println(err)
		}
		log.Println(res)

	}
}
