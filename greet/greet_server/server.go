package main

import (
	"context"
	"fmt"
	"io"
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
	firstName := req.GetGreeting().GetFirstName()
	for i := 0; i < 10; i++ {
		res := &greet.GreetManyTimesResponse {
			Result: "Hello " + firstName,
		}
		stream.Send(res)
		time.Sleep(100 * time.Millisecond)
	}
	log.Printf("GreetManyTimes function completed")
	return nil
}

//LongGreet(GreetService_LongGreetServer) error
func (*greetServer) LongGreet(stream greet.GreetService_LongGreetServer) error {
	log.Printf("Long greet was invoked")

	result := ""
	for {
		data, err := stream.Recv()

		if err == io.EOF {
			return stream.SendAndClose(&greet.LongGreetResponse{Result: result})
		}

		if err != nil {
			log.Fatalf("Failed to create Long Greet %v", err)
			return err
		}

		firstName := data.GetGreeting().GetFirstName()
		if firstName != "" {
			log.Printf("Got Data %v", firstName)
			time.Sleep(100 * time.Millisecond)
			result += firstName
		}
	}
}


func (*greetServer) GreetEveryone(stream greet.GreetService_GreetEveryoneServer) error {
	log.Printf("GreetEveryone was invoked")

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			fmt.Println("Got EOF")
			return nil
		}

		if err != nil {
			log.Fatalf("GreetEveryone while reading client stream %v", err)
			return err
		}

		firstName := req.GetGreeting().GetFirstName()
		result := "Hello " + firstName

		if firstName != "" {
			err = stream.Send(&greet.GreetEveryoneResponse{Result: result})
			fmt.Println("Sending message to " + firstName)
			if err != nil {
				log.Fatalf("GreetEveryone while sending %v", err)
				return err
			}
		}
	}

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
