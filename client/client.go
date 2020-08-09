package main

import (
	"context"
	"fmt"
	"log"

	greet_proto "example.com/greet/greetpb/greet.proto"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello I am a client")
	cc, err := grpc.Dial("0.0.0.0:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Failed to make a connection to server %v", err)
	}

	defer cc.Close()

	c := greet_proto.NewGreetServiceClient(cc)

	doUnary(c)
}

func doUnary(c greet_proto.GreetServiceClient) {
	resp, err := c.Greet(context.Background(), &greet_proto.GreetRequest{
		Greeting: &greet_proto.Greeting{
			FirstName: "Alex",
			LastName:  "Fallenstedt",
		},
	})
	if err != nil {
		log.Fatalf("Failed while calling Greet RPC %v", err)
	}

	log.Printf("Response %v", resp)
}
