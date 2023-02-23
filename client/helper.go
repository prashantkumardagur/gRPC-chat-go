package main

import (
	"fmt"
	"io"

	pb "github.com/prashantkumardagur/grpc-chat/chat"
)

//==============================================================================

func Auth(client pb.ChatClient, stream pb.Chat_MessagingClient) (string, error) {

	// get username from user
	var username string
	fmt.Print("BOT> Enter your username: ")
	Input(&username)

	// send first message to server to register the user
	err := stream.Send(&pb.Message{From: username, To: "/BOT", Text: "password"})
	HandleError(err)

	// return the username after registering the user
	return username, nil
}

//==============================================================================

func Reciever(stream pb.Chat_MessagingClient, waitc chan int, username string) {
	for {
		// recieve messages from the server
		msg, err := stream.Recv()

		// if the user has disconnected
		if err == io.EOF {
			break
		}
		HandleError(err)

		// print the message
		fmt.Print("\r" + msg.GetFrom() + "> " + msg.GetText() + "\n" + username + "> ")
	}
	waitc <- 1
}
