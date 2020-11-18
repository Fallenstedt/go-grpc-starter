package main

import (
	"context"
	"example.com/greet/gen/greet/proto"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"sync"
	"time"
)

func main() {
	fmt.Println("Hello I am a client")
	cc, err := grpc.Dial("0.0.0.0:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Failed to make a connection to server %v", err)
	}

	defer cc.Close()

	c := greet.NewGreetServiceClient(cc)

	//time.Sleep(100 * time.Millisecond)
	//doUnaryGreet(c)
	//time.Sleep(100 * time.Millisecond)
	//doServerStream(c)
	//time.Sleep(100 * time.Millisecond)
	//doClientStream(c)
	//time.Sleep(100 * time.Millisecond)
	doBiStreaming(c)
}

func doBiStreaming(c greet.GreetServiceClient) {
	fmt.Println("Bidirectional streaming")

	stream, err := c.GreetEveryone(context.Background())

	if err != nil {
		log.Fatalf("Error while calling GreetEveryone", err)
	}

	msgs := make([]greet.GreetEveryoneRequest, 10)
	for i := 0; i < 10; i++ {
		msgs = append(msgs, greet.GreetEveryoneRequest{Greeting: &greet.Greeting{
			FirstName: "Alex",
			LastName: "Fallenstedt",
		}})
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, req := range msgs {
			if req.GetGreeting().GetFirstName() != "" {
				stream.Send(&req)
			}
			time.Sleep(time.Millisecond * 100)
		}
		stream.CloseSend()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			res, err := stream.Recv()

			if err == io.EOF {
				fmt.Println("EOF!")
				break
			}

			if err != nil {
				log.Fatalf("GreetEveryone got error! %v", err)
				break
			}

			fmt.Println(res.GetResult())
		}
	}()

	wg.Wait()
	fmt.Println("done")
}

func doClientStream(c greet.GreetServiceClient) {
	stream, err := c.LongGreet(context.Background())

	if err != nil {
		log.Fatalf("Error while calling LongGreet", err)
	}

	msgs := make([]greet.LongGreetRequest, 10)
	for i := 0; i < 10; i++ {
		msgs = append(msgs, greet.LongGreetRequest{Greeting: &greet.Greeting{
			FirstName: "Alex",
			LastName: "Fallenstedt",
		}})
	}

	stream.Send(&greet.LongGreetRequest{Greeting: &greet.Greeting{
		FirstName: "Alex",
		LastName: "Fallenstedt",
	}})
	for _, req := range msgs {
		stream.Send(&req)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Error while calling LongGreet", err)
	}

	log.Printf("Client streaming response %v", res)

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
		log.Printf("Server stream response %v", msg)
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

	log.Printf("Unary Response %v", resp)
}
