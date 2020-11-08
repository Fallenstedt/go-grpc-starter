package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"example.com/greet/gen/greet/proto"
	"google.golang.org/grpc"
)

type greetServer struct{
	greet.GreetServiceServer
}

func (*greetServer) Greet(ctx context.Context, req *greet.GreetRequest) (*greet.GreetResponse, error) {
	log.Printf("Greet function was invoked with %v", req)
	firstName := req.GetGreeting().GetFirstName()
	result := "Hello " + firstName
	res := &greet.GreetResponse{
		Result: result,
	}
	return res, nil
}

func (*greetServer) GreetManyTimes(req *greet.GreetManyTimesRequest, stream greet.GreetService_GreetManyTimesServer) error {
	log.Printf("GreetManyTimes function was invoked with %v", req)
	firstName := req.GetGreeting().GetFirstName();
	for i := 0; i < 10; i++ {
		res := &greet.GreetManyTimesResponse {
			Result: "Hello " + firstName,
		}
		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}
	log.Printf("GreetManyTimes function completed")
	return nil
}

func main() {
	fmt.Println("Hello World")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	greet.RegisterGreetServiceServer(s, &greetServer{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
