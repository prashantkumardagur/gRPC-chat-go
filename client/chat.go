package main

import (
	"context"
	"fmt"
	"io"

	pb "github.com/prashantkumardagur/grpc-chat/chat"
)

//==============================================================================

func Chat(client pb.ChatClient, stream pb.Chat_MessagingClient, username string) {
	var friend string = ""
	var message string = ""

	err := stream.Send(&pb.Message{From: username, To: friend, Text: "blank"})
	HandleError(err)

	// ============================================================================

	waitc := make(chan struct{})

	go func() {
		for {
			msg, err := stream.Recv()
			if err == io.EOF {
				break
			}
			HandleError(err)
			fmt.Print("\r" + msg.GetFrom() + "> " + msg.GetText() + "\n" + username + "> ")
		}
		close(waitc)
	}()

	// ============================================================================

	for {
		fmt.Print(username + "> ")
		Input(&message)

		switch message {

		case "/logout":
			{
				stream.CloseSend()
				<-waitc
				return
			}

		case "/users":
			{
				res, err := client.GetUsers(context.Background(), &pb.Empty{})
				HandleError(err)
				for {
					user, err := res.Recv()
					if err != nil {
						break
					}
					fmt.Println("BOT> " + user.GetUsername())
				}
			}

		case "/connect":
			{
				var temp string
				fmt.Print("BOT> Enter username: ")
				Input(&temp)
				if temp == username {
					fmt.Println("BOT> You can't connect to yourself")
					continue
				}
				res, err := client.CheckUser(context.Background(), &pb.User{Username: temp})
				HandleError(err)
				if res.GetSuccess() {
					fmt.Println("BOT> Connected to " + temp)
					friend = temp
				} else {
					fmt.Println("BOT> User not found")
				}
			}

		case "/broadcast":
			{
				fmt.Println("BOT> Connected to broadcast")
				friend = "/broadcast"
			}

		default:
			{
				if friend == "" {
					fmt.Println("BOT> Please connect to a user")
					continue
				}
				err := stream.Send(&pb.Message{From: username, To: friend, Text: message})
				HandleError(err)
			}
		}
	}
}
