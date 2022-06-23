package main

import (
	"context"
	"flag"
	"io"
	"log"
	"time"

	pb "examples.com/grpc-test/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

var (
	serverAddr = flag.String("addr", "localhost:50051", "The server address in the format of host:port")
)

func downloadEvents(client pb.IotServerClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Hour)
	defer cancel()
	stream, err := client.GetEvents(ctx, &emptypb.Empty{})
	if err != nil {
		log.Fatalf("client.GetEvents failed during connectio: %v", err)
	}
	for {
		event, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("client.GetEvents failed on receiving: %v", err)
		}
		log.Printf("Status: %s", event.Status)
	}
}

func main() {
	flag.Parse()
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial(*serverAddr, opts...)

	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewIotServerClient(conn)
	downloadEvents(client)
}
