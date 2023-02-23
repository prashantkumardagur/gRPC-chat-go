package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

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

var Users = make(map[string]chan *pb.Message)

var collection *mongo.Collection

//==============================================================================

func main() {
	// setting up mongoDB connection
	clientOptions := options.Client().ApplyURI("mongodb+srv://prashantkumar:Password024680@testcluster.8xzqf.mongodb.net/ecommerce?retryWrites=true&w=majority")

	// Connect to MongoDB
	client, dberr := mongo.Connect(context.Background(), clientOptions)
	HandleError(dberr)

	// get collection from database
	collection = client.Database("grpc-chat").Collection("users")

	//==========================================================================

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
