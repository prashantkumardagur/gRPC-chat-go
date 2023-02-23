package main

import (
	"io"

	pb "github.com/prashantkumardagur/grpc-chat/chat"
)

//===============================================================================================================

func Authenticator(stream pb.Chat_MessagingServer) (string, error) {
	// gets first message from the user to get the username and register the user
	message, err := stream.Recv()
	HandleError(err)

	// extract the username from the message
	username := message.GetFrom()
	broadcast(&pb.Message{From: "ChatBOT", To: "/broadcast", Text: username + " has joined the chat"})
	Users[username] = make(chan *pb.Message)

	// return the username after registering the user
	return username, nil
}

//===============================================================================================================

func (s *server) Messaging(stream pb.Chat_MessagingServer) error {

	// gets first message from the user to get the username and register the user
	username, autherr := Authenticator(stream)
	HandleError(autherr)

	// goroutine to handle the user's messages
	go handleChat(stream, Users[username])

	//===========================================================================

	// loop to handle messages from the user
	for {
		message, err := stream.Recv()

		// if the user has disconnected
		if err == io.EOF {
			Users[username] <- &pb.Message{From: "/BOT/", To: username, Text: "Disconnected"}
			delete(Users, username)
			broadcast(&pb.Message{From: "ChatBOT", To: "/broadcast", Text: username + " has left the chat"})
			break
		}
		HandleError(err)

		// check if the message is a broadcast message or a private message
		if message.GetTo() == "/broadcast" {
			broadcast(message)
			continue
		}

		// check if the user is online
		if _, ok := Users[message.GetTo()]; ok {
			Users[message.GetTo()] <- message

			// if the user is not online
		} else {
			Users[username] <- &pb.Message{From: "BOT", To: "", Text: "User " + message.GetTo() + " logged out. Please change chennal."}
		}

	}
	return nil
}

//===============================================================================================================

// handleChat sends a message to the user
func handleChat(stream pb.Chat_MessagingServer, message chan *pb.Message) {
	var msg *pb.Message
	for {
		msg = <-message
		if msg.GetFrom() == "/BOT/" {
			break
		}
		err := stream.Send(msg)
		HandleError(err)
	}
}

//===============================================================================================================

// broadcast sends a message to all users except the sender
func broadcast(message *pb.Message) {
	for username, userChan := range Users {
		if username == message.GetFrom() {
			continue
		}
		userChan <- message
	}
}
