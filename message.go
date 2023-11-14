package main

import (
	"context"
	"io"
	"log"
	"sync"
	// "io"
)

var (
	wg sync.WaitGroup
)

func StreamFromClient(client MutualServiceClient, message string) {
	// Create a stream by calling the streaming RPC method
	// log.Printf(message)
	stream, err := client.StreamFromClient(context.Background())
	if err != nil {
		log.Fatalf("could not get stream: %v", err)
	}

	// Send a message to the server using the stream
	err = stream.Send(&Message{Content: message, SenderId: "client" + addr})
	if err != nil {
		log.Fatalf("could not send message to server: %v", err)
	}

}

func (s *MutualService) StreamFromClient(msgStream MutualService_StreamFromClientServer) error {
	for {
		// get the next message from the stream
		msg, err := msgStream.Recv()

		// the stream is closed so we can exit the loop
		if err == io.EOF {
			break
		}
		// some other error
		if err != nil {
			return err
		}
		// log the message
		log.Printf("Received message from %s: %s", msg.SenderId, msg.Content)
	}
	return nil
}