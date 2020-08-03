package main

import (
	"fmt"
	"log"
	"net"

	greet_proto "example.com/greet/greetpb/greet.proto"
	"google.golang.org/grpc"
)

type server struct{}

func main() {
	fmt.Println("Hello World")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	greet_proto.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
