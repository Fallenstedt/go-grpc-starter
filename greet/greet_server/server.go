package main

import (
	"time"
	"context"
	"fmt"
	"log"
	"net"

	greet_proto "example.com/greet/greetpb/greet.proto"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Greet(ctx context.Context, req *greet_proto.GreetRequest) (*greet_proto.GreetResponse, error) {
	log.Printf("Greet function was invoked with %v", req)
	firstName := req.GetGreeting().GetFirstName()
	result := "Hello " + firstName
	res := &greet_proto.GreetResponse{
		Result: result,
	}
	return res, nil
}

func (*server) GreetManyTimes(req *greet_proto.GreetManyTimesRequest, stream greet_proto.GreetService_GreetManyTimesServer) (error) {
	log.Printf("GreetManyTimes function was invoked with %v", req)
	firstName := req.GetGreeting().GetFirstName();
	for i := 0; i < 10; i++ {
		res := &greet_proto.GreetManyTimesResponse {
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
	greet_proto.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
