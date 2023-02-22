package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	pb "github.com/prashantkumardagur/grpc-chat/chat"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

//==============================================================================

func HandleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Input(str *string) {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\r\n", "", 1)
	*str = text
}

//==============================================================================

func main() {
	// Create a gRPC connection with the server
	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	HandleError(err)
	defer conn.Close()

	// Create a gRPC client
	client := pb.NewChatClient(conn)

	// ============================================================================

	fmt.Println("BOT> Welcome to gRPC Chat")
	var username string

	for {
		fmt.Print("BOT> Enter new username: ")
		Input(&username)

		res, err := client.CheckUser(context.Background(), &pb.User{Username: username})
		HandleError(err)

		if res.GetSuccess() {
			fmt.Println("BOT> User already exists")
		} else {
			fmt.Println("BOT> User logged in")
			break
		}
	}

	stream, err := client.Messaging(context.Background())
	HandleError(err)
	Chat(client, stream, username)

}

// ============================================================================

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
