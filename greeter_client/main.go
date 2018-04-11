package main

import (
	"log"
	"os"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	pb "grpc.io/tutorial/helloworld/helloworld"
)

const (
	// for talking to docker container locally
	// address = "grpc-docker.example.com:50051"
	// for testing ingress
	address = "grpc-ingress.example.com:443"
	//address = "localhost:50051"
	defaultName = "world"
)

func main() {
	creds, err := credentials.NewClientTLSFromFile("certs/tls.crt", "")
	if err != nil {
		log.Fatalf("could not load tls cert: %s", err)
	}

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// contact server, print out response
	name := defaultName

	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: name})

	if err != nil {
		log.Fatalf("could not greet: %+v", err)
	}

	log.Printf("Greeting: %s", r.Message)

	r, err = c.SayHelloAgain(context.Background(), &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Printf("Greeting: %s", r.Message)
}
