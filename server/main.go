package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/prashantkumardagur/grpc-chat/chat"
)

//==============================================================================

func HandleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type server struct {
	pb.ChatServer
}

var Users = make(map[string]chan pb.Message)

//==============================================================================

func main() {
	// Create a listener on TCP port 8080
	lis, netErr := net.Listen("tcp", ":8080")
	HandleError(netErr)

	// Create a gRPC server object
	grpcServer := grpc.NewServer()

	// Register the chat service with the gRPC server
	pb.RegisterChatServer(grpcServer, &server{})

	// Attach gRPC server to the listener
	err := grpcServer.Serve(lis)
	HandleError(err)
}
