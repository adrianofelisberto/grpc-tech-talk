package services

import (
	"context"
	"fmt"
	"time"
	"log"
	"io"
	"github.com/adrianofelisberto/grpc-tech-talk/pb"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
}

func NewUserService() *UserService {
	return &UserService{}
}

func (*UserService) AddUser(ctx context.Context, req *pb.User) (*pb.User, error) {

	fmt.Println("========= AddUser ==========")

	fmt.Println("User -> ", req.Name)

	return &pb.User{
		Id:    req.GetId(),
		Name:  req.GetName(),
	}, nil
}

func (*UserService) AddUserReturnStream(req *pb.User, stream pb.UserService_AddUserReturnStreamServer) error {

	fmt.Println("========= AddUserReturnStream ==========")

	fmt.Println("User sending stream -> ", req.Name)

	fmt.Println("Init")
	stream.Send(&pb.UserResultStream{
		Status: "Init",
		User: &pb.User{},
	})

	time.Sleep(time.Second * 3)

	fmt.Println("Insert")
	stream.Send(&pb.UserResultStream{
		Status: "Insert",
		User: &pb.User{},
	})

	time.Sleep(time.Second * 3)

	fmt.Println("Finish")
	stream.Send(&pb.UserResultStream{
		Status: "Finish",
		User: &pb.User{
			Id: req.GetId(),
			Name: req.GetName(),
		},
	})

	time.Sleep(time.Second * 3)

	return nil
}

func (*UserService) AddUsers(stream pb.UserService_AddUsersServer) error {

	fmt.Println("========= AddUsers ==========")
	users := []*pb.User{}

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.Users{
				User: users,
			})
		}
		if err != nil {
			log.Fatalf("Error receiving stream: %v", err)
		}

		fmt.Println("User receiving stream -> ", req.Name)

		users = append(users, &pb.User{
			Id:    req.GetId(),
			Name:  req.GetName(),
		})
		fmt.Println("Adding", req.GetName())

	}

}

func (*UserService) AddUserStream(stream pb.UserService_AddUserStreamServer) error {

	fmt.Println("========= AddUserStream ==========")
	for {
		req, err := stream.Recv()
		
		
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error receiving stream from the client: %v", err)
		}
		
		fmt.Println("User receiving stream -> ", req.Name)

		fmt.Println("User sending stream -> ", req.Name)
		
		err = stream.Send(&pb.UserResultStream{
			Status: "Added",
			User:   req,
		})

		if err != nil {
			log.Fatalf("Error sending stream to the client: %v", err)
		}
	}
}