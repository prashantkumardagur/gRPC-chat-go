package main

import (
	"context"
	"io"
	"log"

	pb "github.com/prashantkumardagur/grpc-chat/chat"
)

//==============================================================================

func (s *server) GetUsers(req *pb.Empty, stream pb.Chat_GetUsersServer) error {
	for username := range Users {
		stream.Send(&pb.User{Username: username})
	}
	return nil
}

//==============================================================================

func (s *server) CheckUser(ctx context.Context, req *pb.User) (*pb.Response, error) {
	username := req.GetUsername()
	if _, ok := Users[username]; ok {
		return &pb.Response{Message: "User found", Success: true}, nil
	}
	return &pb.Response{Message: "User not found", Success: false}, nil
}

//==============================================================================

func (s *server) Messaging(stream pb.Chat_MessagingServer) error {

	message, err := stream.Recv()
	HandleError(err)
	var username = message.GetFrom()
	broadcast(pb.Message{From: "ChatBOT", To: "/broadcast", Text: username + " has joined the chat"})
	Users[username] = make(chan pb.Message)
	log.Println("User " + username + " connected")

	go handleChat(stream, Users[username])

	for {
		message, err := stream.Recv()
		if err == io.EOF {
			log.Print("User " + username + " disconnected")
			Users[username] <- pb.Message{From: "/BOT/", To: "", Text: "You have been disconnected"}
			delete(Users, username)
			broadcast(pb.Message{From: "ChatBOT", To: "/broadcast", Text: username + " has left the chat"})
			break
		}
		HandleError(err)
		if message.GetTo() == "/broadcast" {
			broadcast(*message)
		} else {
			Users[message.GetTo()] <- *message
		}
	}
	return nil
}

func handleChat(stream pb.Chat_MessagingServer, message chan pb.Message) {
	for msg := range message {
		if msg.GetFrom() == "/BOT/" {
			break
		}
		stream.Send(&msg)
	}
}

func broadcast(message pb.Message) {
	for username, userChan := range Users {
		if username == message.GetFrom() {
			continue
		}
		userChan <- message
	}
}
