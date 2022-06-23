package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	pb "examples.com/grpc-test/protos"
	"google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

var (
	port     = flag.Int("port", 50051, "The server port")
	statuses = []string{`Starting`, `Online`, `KO`}
)

type server struct {
	pb.UnimplementedIotServerServer
}

func (c *server) GetEvents(empty *emptypb.Empty, stream pb.IotServer_GetEventsServer) error {
	for {
		if err := stream.Send(&pb.Event{Status: statuses[rand.Intn(len(statuses))]}); err != nil {
			return err
		}
		time.Sleep(time.Second)
	}
	return nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterIotServerServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
