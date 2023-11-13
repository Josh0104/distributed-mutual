package main

import (
	"context"
	"log"
)



func StreamFromClient(client MutualServiceClient, port_ string) {
	log.Printf("Client %s is streaming to server", port_)
	// Create a stream by calling the streaming RPC method
	log.Printf("Client %s is streaming to server", addr)
	stream, err := client.StreamFromClient(context.Background())
	if err != nil {
		log.Fatalf("could not get stream: %v", err)
	}

	// Send a message to the server using the stream
	err = stream.Send(&Message{Content: "Hello from client", SenderId: "client" + addr})
	if err != nil {
		log.Fatalf("could not send message to server: %v", err)
	}

	// ... Continue interacting with the stream as needed.
	// You can send more messages or receive responses from the server.
}