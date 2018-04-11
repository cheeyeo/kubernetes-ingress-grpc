package main

import (
	"log"
	//"crypto/tls"
	"net"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	pb "grpc.io/tutorial/helloworld/helloworld"
)

const (
	port = ":50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	glog.Infof("Inside SayHello: %v\n", in)
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func (s *server) SayHelloAgain(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	glog.Infof("Inside SayHelloAgain: %v\n", in)
	return &pb.HelloReply{Message: "Hello again " + in.Name}, nil
}

func main() {
	// Recommended from nginx blog
	// Sets grpcs on the tcp
	// cer, err := tls.LoadX509KeyPair("certs/tls.crt", "certs/tls.key")
	// if err != nil {
	//   glog.Fatal("Failed to load key pair: %v\n", err)
	// }

	// config := &tls.Config{ Certificates: []tls.Certificate{cer} }
	// lis, err := tls.Listen("tcp", port, config)
	// if err != nil {
	//   glog.Fatal("Failed to listen: %v\n", err)
	// }

	lis, err := net.Listen("tcp", port)
	if err != nil {
		glog.Fatal("Failed to listen: %v\n", err)
	}

	// grab the tls certs
	creds, err := credentials.NewServerTLSFromFile("certs/tls.crt", "certs/tls.key")
	if err != nil {
		glog.Fatalf("Could not load TLS certs %s", err)
	}

	//create array of grpc options with creds
	opts := []grpc.ServerOption{
		grpc.MaxConcurrentStreams(200),
		grpc.UnaryInterceptor(loggingUnaryInterceptor),
		grpc.Creds(creds),
	}

	s := grpc.NewServer(opts...)
	pb.RegisterGreeterServer(s, &server{})

	// Test reflection with grpcurl
	reflection.Register(s)

	glog.Infof("Greeter Server running on TCP port %s", port)

	if err := s.Serve(lis); err != nil {
		glog.Fatal("Failed to serve: %v\n", err)
		log.Fatalf("failed to serve: %v", err)
	}
}

func loggingUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	glog.Infof("Unary Request %s", info.FullMethod)
	return handler(ctx, req)
}
