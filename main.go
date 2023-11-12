package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	sync "sync"
	"time"

	"github.com/manifoldco/promptui"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type MutualService struct {
	UnimplementedMutualServiceServer
}

var (
	port_numbers = [3]string{":5001", ":5002", ":5003"}
	otherClients = [2]string{}
	clientFlag   = flag.Int("client", 0, "Enter client number 1, 2 or 3")
	addr         string
	wgServer     sync.WaitGroup
	wgClients    sync.WaitGroup
	globalCancel context.CancelFunc // Declare at a higher scope
	clients_server       MutualServiceClient
)

func main() {
	flag.Parse()
	log.Printf("Listening on %s", addr)
	// startServer()
	// enterCriticalSection(nameFlag)
	configClient()
	wgServer.Wait()
	wgClients.Wait()
	log.Default().Println("Done")
}

func configClient() {

	//switch case
	switch *clientFlag {
	case 1:
		addr = port_numbers[0]
		otherClients = [2]string{"localhost:5002", "localhost:5003"}
	case 2:
		addr = port_numbers[1]
		otherClients = [2]string{port_numbers[0], port_numbers[2]}
	case 3:
		addr = port_numbers[2]
		otherClients = [2]string{port_numbers[0], port_numbers[1]}
	default:
		log.Fatal("Please enter a valid client number with the command line flag -client 1, 2 or 3}")
	}
	// Increment the WaitGroup counter for each goroutine
	wgServer.Add(1)

	go startServer()
	go promptForInput()

}

func startServer() {
	defer wgServer.Done()
	log.Printf("Trying to start server on %s", addr)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("cannot create listener: %s", err)
		return
	}
	serverRegistrar := grpc.NewServer()
	service := &MutualService{}
	RegisterMutualServiceServer(serverRegistrar, service)
	err = serverRegistrar.Serve(lis)
	if err != nil {
		log.Fatalf("impossible to serve: %s", err)
		return
	}
	log.Printf("Successful, starting server on %s", addr)
}

func promptForInput() {
	active := true

	for active {
		prompt := promptui.Select{
			Label:   "Select an option",
			Items: []string{"Connect to other clients", "Enter critical section", "Exit"},
		}
		_, input, err := prompt.Run()
		if err != nil {
			log.Fatalf("could not prompt: %v", err)
		}

		if input == "Connect to other clients" {
			active = false
			connectToOtherClients()
		} else if input == "Enter critical section" {
			active = false
			enterCriticalSection(addr)
		} else if input == "Exit" {
			active = false
			exit()
		} else {
			log.Fatalf("Invalid input")
		}

	}
}

func enterCriticalSection(name string) {
	defer leaveCriticalSection(name)
	log.Printf("Client %s entering critical section", name)
}

func leaveCriticalSection(name string) {
	log.Printf("Client %s leaving critical section", name)
}

func connectToOtherClients() {
	log.Printf("Connecting to other clients")
	wgClients.Add(2)
	for _, port := range otherClients {
		defer wgClients.Done()
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		//dial options
		//the server is not using TLS, so we use insecure credentials
		//(should be fine for local testing but not in the real world)
		opts := []grpc.DialOption{
			grpc.WithBlock(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		}
		log.Printf("Dialing %s", port)
		conn, err := grpc.DialContext(ctx, fmt.Sprintf("%s", addr), opts...)
		if err != nil {
			log.Printf("Fail to Dial : %v", err)
			return
		}
		res, err := clients_server.Join(ctx, &JoinRequest{
			NodeName: addr,
		})
		if err != nil {
			log.Printf("Fail to Join : %v", err)
			return
		}
		log.Printf("Join response: %v", res)

		log.Println("the connection is: ", conn.GetState().String())

		log.Printf("Connected to %s", port)
	}
}

func exit() {
	log.Printf("Exiting")
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
}
