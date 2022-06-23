package main

import (
	"context"
	pb "examples.com/grpc-test/protos"
	"flag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"time"
)

var (
	serverAddr = flag.String("addr", "localhost:50051", "The server address in the format of host:port")
)

func main() {

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial(*serverAddr, opts...)

	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	defer conn.Close()

	client := pb.NewGameClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Hour)
	defer cancel()

	stream, err := client.Ping(ctx)
	if err != nil {
		log.Fatalf("client.GetEvents failed during connectio: %v", err)
	}
	for {
		log.Printf(`Sending a ping`)
		stream.Send(&pb.Status{})
		pong, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("client.GetEvents failed on receiving: %v", err)
		}
		log.Printf("Rceived a %s", pong.Status)
	}
	log.Printf(`Done`)
}
