package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"

	"github.com/yurianxdev/grpc-course/greet/greetpb"
)

type server struct{}

func (s *server) GreetEveryOne(stream greetpb.GreetService_GreetEveryOneServer) error {
	log.Printf("Request for GreetEveryOne was accepted...\n")

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		firstName := req.GetGreeting().GetFirstName()
		result := "Hello " + firstName + "! "
		err = stream.Send(&greetpb.GreetEveryOneResponse{
			Response: result,
		})
		if err != nil {
			return err
		}
	}
}

func (s *server) Greet(_ context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	log.Println("Request for Greet was accepted...")
	result := "Hello " + req.Greeting.FirstName
	res := &greetpb.GreetResponse{
		Result: result,
	}

	return res, nil
}

func (s *server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	log.Printf("Request for GreetManyTimes was accepted...\n")
	firstName := req.GetGreeting().FirstName
	for i := 1; i <= 10; i++ {
		result := fmt.Sprintf("Hello %s, this is the greet number %d", firstName, i)
		res := &greetpb.GreetManyTimesResponse{
			Result: result,
		}
		err := stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *server) LongGreet(greetServer greetpb.GreetService_LongGreetServer) error {
	log.Printf("Request for LongGreet was accepted...\n")
	result := "Hello "

	for {
		req, err := greetServer.Recv()
		if err == io.EOF {
			return greetServer.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}

		firstName := req.GetGreeting().GetFirstName()
		result += firstName + " "
	}
}

func main() {
	port := "50051"
	log.Println("Starting server...")
	li, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Error listening server: %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	log.Printf("Server listening on port %s", port)
	if err := s.Serve(li); err != nil {
		log.Fatalf("Error serving: %v", err)
	}
}
