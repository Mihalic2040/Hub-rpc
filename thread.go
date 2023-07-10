package hubrpc

import (
	"fmt"

	"github.com/Mihalic2040/Hub-rpc.git/src/proto/api"
)

func Thread(handlers HandlerMap, data *api.Request) (api.Response, error) {

	// Call a specific handler by name
	handlerName := data.Handler
	handler, ok := handlers[handlerName]
	if !ok {
		//fmt.Printf("Handler '%s' not found\n", handlerName)
		return api.Response{
			Payload: "Handler not found",
			Status:  500,
		}, nil
	}

	// Call the handler function with the input data
	output, err := handler(data)
	//handler(inputData)
	if err != nil {
		fmt.Printf("[SERVER: Thread] Error executing handler: %v\n", err)
		return api.Response{
			Payload: err.Error(),
			Status:  500,
		}, nil
	}

	return output, nil

}
