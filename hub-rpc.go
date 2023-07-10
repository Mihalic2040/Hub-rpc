package hubrpc

import (
	"bufio"
	"fmt"
	"log"

	"github.com/Mihalic2040/Hub-rpc.git/src/proto/api"
	"github.com/libp2p/go-libp2p/core/network"
	"google.golang.org/protobuf/proto"
)

type Handler func(input *api.Request) (response api.Response, err error)

// HandlerMap holds a map of handler names to their corresponding functions
type HandlerMap map[string]Handler

func (m *HandlerMap) HandleFunc(handlerName string, handlerFunc Handler) {
	(*m)[handlerName] = handlerFunc
}

func Stream_handler(stream network.Stream, handlers HandlerMap) {
	//log.Println("New stream!!")
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	// Create a data channel to receive the information
	dataCh := make(chan []byte)

	// Goroutine to read information from the stream and send it to the data channel
	go func() {
		for {
			// Read a chunk of data from the stream
			data := make([]byte, 1024) // Adjust the buffer size as per your needs
			n, err := rw.Read(data)
			if err != nil {
				//log.Println("Error reading data from stream:", err)
				close(dataCh)
				break
			}

			// Send the received data to the data channel
			dataCh <- data[:n]
		}
	}()

	// Read the information from the data channel
	data := <-dataCh

	// Create a new instance of your protobuf message
	req := &api.Request{}

	// Decode the protobuf data
	if err := proto.Unmarshal(data, req); err != nil {
		log.Println("[SERVER] Error decoding protobuf data:", err)
		return
	}

	// Start data processing thread

	response, err := Thread(handlers, req)

	// Send the response back to the client

	// Send the response to the stream
	response_b, err := proto.Marshal(&response)
	if err != nil {
		log.Println("[SERVER] Error encoding protobuf data:", err)

	}

	// Write the request bytes to the stream
	if _, err := rw.Write(response_b); err != nil {
		fmt.Println("[SERVER] Error writing protobuf response:", err)

	}

	// Flush the writer to ensure the data is sent
	if err := rw.Flush(); err != nil {
		log.Println("[SERVER] Error flush cahanel:", err)

	}

	// Close the stream
	stream.Close()

	log.Println("[SERVER] ", "Protocol:", stream.Conn().ConnState().Transport, "Remote:", stream.Conn().RemotePeer(), "Status:", response.Status)

}
