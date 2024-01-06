package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"

	train "path-to-proto/train_ticket"
)

func main() {
	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()

	client := train.NewTrainTicketServiceClient(conn)

	// Implement client calls to gRPC methods here
}
