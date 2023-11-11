package main

import (
	"context"
	"fmt"
	"net"

	pb "github.com/Josh0104/distributed-mutual/proto"
)

type myInvoicerServer struct {
	pb.UnimplementedInvoicerServer
}

func (s myInvoicerServer) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	return &pb.CreateResponse{
		Pdf:  []byte(req.From),
		Docx: []byte("test"),
	}, nil
}
func main() {
	fmt.Println("Hello, World!")
	lis, err := net.Listen("tcp", ":8089")
	if err != nil {
		log.Fatalf("cannot create listener: %s", err)
	}
	serverRegistrar := grpc.NewServer()
	service := &myInvoicerServer{}
	pb.RegisterInvoicerServer(serverRegistrar, service)
	err = serverRegistrar.Serve(lis)
	if err != nil {
		log.Fatalf("impossible to serve: %s", err)
	}
}


