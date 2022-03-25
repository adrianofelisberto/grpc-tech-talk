package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/adrianofelisberto/grpc-tech-talk/pb"
	"google.golang.org/grpc"
)

func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Coud not connect to gRPC Server: %v", err)
	}
	defer connection.Close()

	client := pb.NewUserServiceClient(connection)

	AddUser(client)
	time.Sleep(time.Second * 3)
	AddUserReturnStream(client)
	time.Sleep(time.Second * 3)
	AddUsers(client)
	time.Sleep(time.Second * 3)
	AddUserStream(client)

}

func AddUser(client pb.UserServiceClient) {

	fmt.Println("========= AddUser ==========")

	req := &pb.User{
		Id:    "0",
		Name:  "Adriano",
	}

	res, err := client.AddUser(context.Background(), req)

	if err != nil {
		log.Fatalf("Coud not make gRPC request: %v", err)
	}

	fmt.Println(res)
}

func AddUserReturnStream(client pb.UserServiceClient) {

	fmt.Println("========= AddUserReturnStream ==========")

	req := &pb.User{
		Id:    "0",
		Name:  "Adriano",
	}

	responseStream, err := client.AddUserReturnStream(context.Background(), req)

	if err != nil {
		log.Fatalf("Coud not make gRPC request: %v", err)
	}

	for {
		stream, err := responseStream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Coud not receive the msg: %v", err)
		}
		fmt.Println("Status: ", stream.Status)
	}

}

func AddUsers(client pb.UserServiceClient) {
	fmt.Println("========= AddUsers ==========")
	reqs := []*pb.User{
		&pb.User{
			Id:    "a1",
			Name:  "Adriano 1",
		},
		&pb.User{
			Id:    "a2",
			Name:  "Adriano 2",
		},
		&pb.User{
			Id:    "a3",
			Name:  "Adriano 3",
		},
		&pb.User{
			Id:    "a4",
			Name:  "Adriano 4",
		},
		&pb.User{
			Id:    "a5",
			Name:  "Adriano 5",
		},
		&pb.User{
			Id:    "a6",
			Name:  "Adriano 6",
		},
	}

	stream, err := client.AddUsers(context.Background())

	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	for _, req := range reqs {
		stream.Send(req)
		fmt.Println("Sending user: ", req.Name)
		time.Sleep(time.Second * 3)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error receiving response: %v", err)
	}

	fmt.Println(res)
}

func AddUserStream(client pb.UserServiceClient) {

	fmt.Println("========= AddUserStream ==========")

	stream, err := client.AddUserStream(context.Background())
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	reqs := []*pb.User{
		&pb.User{
			Id:    "a1",
			Name:  "Adriano 1",
		},
		&pb.User{
			Id:    "a2",
			Name:  "Adriano 2",
		},
		&pb.User{
			Id:    "a3",
			Name:  "Adriano 3",
		},
		&pb.User{
			Id:    "a4",
			Name:  "Adriano 4",
		},
		&pb.User{
			Id:    "a5",
			Name:  "Adriano 5",
		},
		&pb.User{
			Id:    "a6",
			Name:  "Adriano 6",
		},
	}

	wait := make(chan int)

	go func() {
		for _, req := range reqs {
			fmt.Println("Sending user: ", req.Name)
			stream.Send(req)
			time.Sleep(time.Second * 2)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error receiving data: %v", err)
				break
			}
			fmt.Printf("Receiving user %v with status: %v\n", res.GetUser().GetName(), res.GetStatus())
		}
		close(wait)
	}()

	<-wait

}
