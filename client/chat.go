package main

import (
	"context"
	"fmt"

	pb "github.com/prashantkumardagur/grpc-chat/chat"
)

//==============================================================================

func Chat(client pb.ChatClient, stream pb.Chat_MessagingClient, username string) {
	var friend string = "/broadcast"
	var message string = ""

	// ============================================================================

	waitc := make(chan int)
	go Reciever(stream, waitc, username)

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
				// get username from user
				var temp string
				fmt.Print("BOT> Enter username: ")
				Input(&temp)
				if temp == username {
					fmt.Println("BOT> You can't connect to yourself")
					continue
				}

				// check if the user exists
				res, err := client.CheckUser(context.Background(), &pb.User{Username: temp})
				HandleError(err)
				if res.GetSuccess() { // if user exists
					fmt.Println("BOT> Connected to " + temp)
					friend = temp
				} else { // if user doesn't exist
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
				err := stream.Send(&pb.Message{From: username, To: friend, Text: message})
				HandleError(err)
			}
		}
	}
}
