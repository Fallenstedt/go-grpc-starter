package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"example.com/greet/gen/greet/proto"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello I am a client")
	cc, err := grpc.Dial("0.0.0.0:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Failed to make a connection to server %v", err)
	}

	defer cc.Close()

	c := greet.NewGreetServiceClient(cc)

	doUnaryGreet(c)
	doServerStream(c)
}

func doServerStream(c greet.GreetServiceClient) {
	resStream, err := c.GreetManyTimes(context.Background(), &greet.GreetManyTimesRequest{
		Greeting: &greet.Greeting {
			FirstName: "Alex",
			LastName: "Fallenstedt",
		},
	})
	if err != nil {
		log.Fatalf("Failed while calling GreetManyTimes RPC %v", err)
	}
	
	for {
		msg, err := resStream.Recv()
		// Server ended stream
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Error while reading stream %v", err)
		}
		log.Printf("Response %v", msg)
	}
	

}

func doUnaryGreet(c greet.GreetServiceClient) {
	resp, err := c.Greet(context.Background(), &greet.GreetRequest{
		Greeting: &greet.Greeting{
			FirstName: "Alex",
			LastName:  "Fallenstedt",
		},
	})
	if err != nil {
		log.Fatalf("Failed while calling Greet RPC %v", err)
	}

	log.Printf("Response %v", resp)
}
