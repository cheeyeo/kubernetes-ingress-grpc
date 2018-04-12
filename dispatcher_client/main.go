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
  // localhost
  //address = "grpc-dev.example.com:50051"
  // k8 ingress
  address = "grpc-dispatcher.example.com:443"
  defaultTitle = "Infin1te"
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

  title := defaultTitle
  if len(os.Args) > 1 {
    title = os.Args[1]
  }

  // creating dispatcher client
  d := pb.NewDispatcherClient(conn)
  dr, err := d.TestDispatch(context.Background(), &pb.DispatchRequest{Title: title})

  if err != nil {
    log.Fatalf("Could not dispatch: %+v", err)
  }

  log.Printf("Dispatcher: %s", dr.Message)
}
