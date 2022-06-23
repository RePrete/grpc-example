package main

import (
	pb "examples.com/grpc-test/protos"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"time"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type server struct {
	pb.UnimplementedGameServer
}

func (s server) Ping(stream pb.Game_PingServer) error {
	for {
		_, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		log.Printf("Received a ping")
		stream.Send(&pb.Pong{Status: `pong`})
		log.Printf("Responding with pong")
		time.Sleep(time.Second)
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGameServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
