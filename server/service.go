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

	// gets first message from the user to get the username and register the user
	message, err := stream.Recv()
	HandleError(err)
	var username = message.GetFrom()
	broadcast(pb.Message{From: "ChatBOT", To: "/broadcast", Text: username + " has joined the chat"})
	Users[username] = make(chan pb.Message)
	log.Println("User " + username + " connected")

	// goroutine to handle the user's messages
	go handleChat(stream, Users[username])

	// loop to handle messages from the user
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
			if _, ok := Users[message.GetTo()]; ok {
				Users[message.GetTo()] <- *message
			} else {
				Users[username] <- pb.Message{From: "BOT", To: "", Text: "User " + message.GetTo() + " logged out. Please change chennal."}
			}
		}
	}
	return nil
}

// handleChat sends a message to the user
func handleChat(stream pb.Chat_MessagingServer, message chan pb.Message) {
	for msg := range message {
		if msg.GetFrom() == "/BOT/" {
			break
		}
		stream.Send(&msg)
	}
}

// broadcast sends a message to all users except the sender
func broadcast(message pb.Message) {
	for username, userChan := range Users {
		if username == message.GetFrom() {
			continue
		}
		userChan <- message
	}
}
