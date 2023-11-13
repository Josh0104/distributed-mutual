package main

import (
	"context"
	"flag"
	"fmt"

	// "io"
	"log"
	"net"
	"sync"
	"time"

	"github.com/manifoldco/promptui"
	"google.golang.org/grpc"
)

type MutualService struct {
	UnimplementedMutualServiceServer
	name                string // Not required but useful if you want to name your server
	port                string // Not required but useful if your server needs to know what port it's listening to
	isInCriticalSection bool
	clientMu            sync.Mutex
	isWaiting           bool
	nextPeer            string
	prevPeer            string
}

var (
	port_numbers       = [3]string{":5001", ":5002", ":5003"}
	otherClientsPort   = [2]string{}
	otherClientsServer []MutualServiceClient
	clientFlag         = flag.Int("client", 0, "Enter client number 1, 2 or 3")
	addr               string
	wgServer           sync.WaitGroup
	wgClients          sync.WaitGroup
	globalCancel       context.CancelFunc  // Declare at a higher scope
	mutual_client      MutualServiceClient //the server
	mutual_server      MutualServiceServer //the client
	otherClientsServerMutex sync.Mutex

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
		otherClientsPort = [2]string{port_numbers[1], port_numbers[2]}
	case 2:
		addr = port_numbers[1]
		otherClientsPort = [2]string{port_numbers[0], port_numbers[2]}
	case 3:
		addr = port_numbers[2]
		otherClientsPort = [2]string{port_numbers[0], port_numbers[1]}
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
	mutual_server := &MutualService{
		name:                "client" + addr,
		port:                addr,
		isInCriticalSection: true,
		isWaiting:           true,
	}

	RegisterMutualServiceServer(serverRegistrar, mutual_server)
	err = serverRegistrar.Serve(lis)
	if err != nil {
		log.Fatalf("impossible to serve: %s", err)
		return
	}

	setToken(mutual_server, false, false)

	log.Printf("Successful, starting server on %s", addr)
}

func promptForInput() {
	active := true
	connected := false

	for active {

		var prompt = promptui.Select{}
		if !connected {
			prompt = promptui.Select{
				Label: "Select an option",
				Items: []string{"Connect to other clients", "Exit"},
			}
		} else {
			prompt = promptui.Select{
				Label: "Select an option",
				Items: []string{"Enter critical section", "Exit", "Send message to clients"},
			}
		}

		_, input, err := prompt.Run()
		if err != nil {
			log.Fatalf("could not prompt: %v", err)
		}

		if input == "Connect to other clients" {
			// active = false
			connectToOtherClients()
			connected = true
		} else if input == "Enter critical section" {
			// active = false
			enterWaitingState()
		} else if input == "Exit" {
			exit()
			active = false
		} else if input == "Send message to clients" {
			log.Printf("Sending message to other clients")
			if input == "Send message to clients" {
				if active {
					prompt := promptui.Prompt{
						Label: "input your message and send ",
						Validate: func(input string) error {
							if len(input) == 0 || len(input) > 128 {
								return fmt.Errorf("input must be at least 1 character or under 128 characters")
							}
							return nil
						},
					}
					input, err := prompt.Run()
					if err != nil {
						log.Fatalf("could not get input: %v", err)
					}

					// sendMessageToClients(otherClientsServer[0], input)
					// sendMessageToClients(otherClientsServer[1], input)
					send(otherClientsServer[0], input)
					send(otherClientsServer[1], input)
					// Example usage to check isInCriticalSection for the first client

				}
			}
		} else {
			log.Fatalf("Invalid input")
		}

	}
}

func enterWaitingState() {
    log.Printf("Client %s entering waiting state", addr)
    active := true
	waitIndex := 0;
    for active {
        otherClientsServerMutex.Lock()
        token1, err := otherClientsServer[0].GetToken(context.Background(), &Token{})
        token2, err := otherClientsServer[1].GetToken(context.Background(), &Token{})
        otherClientsServerMutex.Unlock()

        if err != nil {
            // handle error
        }
		log.Printf("Token 1: %v", token1.IsCritical)

        if token1.IsCritical || token2.IsCritical {
            active = false
            // Initialize the Token with IsWaiting and IsCritical set to true
			// setToken(, true, true)

            otherClientsServerMutex.Lock()

            enterCriticalSection(addr)
        } else if waitIndex == 5{
			log.Println("Client is waiting for too long, exiting waiting state")
			active = false
		}else {
			waitIndex = waitIndex + 1
			log.Printf("Client %s is waiting", addr)
            time.Sleep(5 * time.Second)
        }
    }
}



func enterCriticalSection(name string) {
	defer leaveCriticalSection(name)
	log.Printf("Client %s entering critical section", addr)

	var message = "Client " + addr + " entering critical section"
	StreamFromClient(otherClientsServer[0], message)
	StreamFromClient(otherClientsServer[1], message)
	getTokenInfo(otherClientsServer[0])
	getTokenInfo(otherClientsServer[1])
	time.Sleep(5 * time.Second)
}

func leaveCriticalSection(name string) {
	log.Printf("Client %s leaving critical section", name)
}

func connectToOtherClients() {
	log.Printf("Connecting to other clients")
	wgClients.Add(2)

	for _, port := range otherClientsPort {
		go func(port string) {
			defer wgClients.Done()

			log.Println("Connecting to ", port)

			// Create a gRPC client to interact with the server
			conn, err := grpc.Dial(fmt.Sprintf("%s", port), grpc.WithInsecure())
			if err != nil {
				log.Printf("Failed to dial %s: %v", port, err)
				return
			}
			// defer conn.Close()

			// Initialize the gRPC client
			clients_server := NewMutualServiceClient(conn)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			// StreamFromClient(clients_server, port)

			otherClientsServer = append(otherClientsServer, clients_server)
			// sayHiToClient(clients_server, port)

			res, err := clients_server.Join(ctx, &JoinRequest{
				NodeName: addr,
			})

			res2, err2 := clients_server.GetToken(ctx, &Token{
				TokenName:  addr + "token",
				IsWaiting:  false,
				IsCritical: false,
			})

			if err != nil || err2 != nil {
				log.Printf("Fail to Join: %v", err)
				return
			}

			log.Printf("Join response: %v", res)

			log.Printf("Token response: %v", res2)

			log.Println("the connection is: ", conn.GetState().String())

			log.Printf("Connected to %s", port)

		}(port)
	}
}

func exit() {
	log.Printf("Exiting")
	for _, _clientServer := range otherClientsServer {
		stream, err := _clientServer.StreamFromClient(context.Background())

		stream.Send(&Message{Content: addr, SenderId: "Bye"})
		if err != nil {
			log.Println(err)
			return
		}
		farewell, err := stream.CloseAndRecv()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("server says: ", farewell)
	}
}

func (s *MutualService) Join(ctx context.Context, req *JoinRequest) (*JoinResponse, error) {
	// log.Printf("[Server: %s time: %d] Received A JOIN req from node %s", s.getName(), s.GetLamport(), req.Lamport.GetNodeId())
	// add the client to the broadcast

	return &JoinResponse{
		NodeId: req.GetNodeName(),
	}, nil
}

func sayHiToClient(server MutualServiceClient, port string) {
	// get a stream to the server
	stream, err := server.StreamFromClient(context.Background())
	if err != nil {
		log.Println(err)
		return
	}
	// send some messages to the server
	stream.Send(&Message{Content: "I am now connected to you", SenderId: port})
}

func sendMessageToClients(server MutualServiceClient, input string) {
	stream, err := server.StreamFromServer(context.Background(), &Message{Content: "Hello from client", SenderId: "client" + addr})

	stream.SendMsg(&Message{Content: input, SenderId: "client" + addr})

	if err != nil {
		log.Println(err)
		return
	}

}

// Server receives message from client
func send(client MutualServiceClient, input string) {
	// Create a stream by calling the streaming RPC method
	StreamFromClient(client, input)
}

func getTokenInfo(client MutualServiceClient) {
	// Call the GetServerInfo RPC
	res, err := client.GetToken(context.Background(), &Token{})
	if err != nil {
		log.Printf("Failed to get server info: %v", err)
		return
	}

	// Access the information
	tokenName := res.GetTokenName()
	criticalStatus := res.GetIsCritical()

	log.Printf("Server Info: Name: %s, criticalStatus: %v", tokenName, criticalStatus)
}

func (s *MutualService) GetToken(ctx context.Context, req *Token) (*Token, error) {
	// Return information about the server
	return &Token{
		TokenName:  "111",
		IsWaiting:  false,
		IsCritical: false,
	}, nil
}

func setToken(s *MutualService, isWaiting bool, isCritical bool) {
	
	s.GetToken(context.Background(), &Token{
		TokenName:  "111",
		IsWaiting:  isWaiting,
		IsCritical: isCritical,
	})

}