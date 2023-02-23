package main

import (
	"context"

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
