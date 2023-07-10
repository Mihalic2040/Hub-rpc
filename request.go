package hubrpc

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"

	"github.com/Mihalic2040/Hub-rpc/src/proto/api"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"google.golang.org/protobuf/proto"
)

func NewRequest(peerID string, data *api.Request, protocolId string, dht dht.IpfsDHT, host host.Host) (*api.Response, error) {
	// Find a peer by its ID
	targetPeerID, err := peer.Decode(peerID)
	if err != nil {
		return nil, fmt.Errorf("Invalid peer ID: %v", err)
	}

	peerInfo, err := dht.FindPeer(context.Background(), targetPeerID)
	if err != nil {
		return nil, fmt.Errorf("Fail to find peer: %v", err)
	}

	// Create a stream to the peer
	stream, err := host.NewStream(context.Background(), peerInfo.ID, protocol.ID(protocolId))
	if err != nil {
		return nil, fmt.Errorf("Failed to create stream: %v", err)
	}

	// Create a bufio ReadWriter using the stream
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	// Create a request

	// Serialize the request to bytes
	bytes, err := proto.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("Failed to serialize request: %v", err)
	}

	// Write the request bytes to the stream
	if _, err := rw.Write(bytes); err != nil {
		return nil, fmt.Errorf("Failed to send request: %v", err)
	}

	// Flush the writer to ensure the data is sent
	if err := rw.Flush(); err != nil {
		return nil, fmt.Errorf("Error flusing writer: %v", err)
	}

	// Read the response from the stream
	responseBytes, err := ioutil.ReadAll(rw)
	if err != nil {
		return nil, fmt.Errorf("Error: %v", err)
	}

	// Create a response message to unmarshal the response bytes
	response := &api.Response{}

	// Unmarshal the response bytes
	if err := proto.Unmarshal(responseBytes, response); err != nil {
		return nil, fmt.Errorf("Error: %v", err)
	}

	//log.Println(response)
	// // Close the stream only if a response is received
	stream.Close()

	// Use the response message as needed
	// ...
	return response, nil
}
