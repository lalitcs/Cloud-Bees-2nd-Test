package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"

	train "path-to-proto/train_ticket"
)

const (
	address = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()

	client := train.NewTrainTicketServiceClient(conn)

	// Example usage of gRPC methods
	purchaseTicketExample(client)
	getReceiptDetailsExample(client)
	viewUserSeatsExample(client)
	removeUserExample(client)
	modifyUserSeatExample(client)
}

func purchaseTicketExample(client train.TrainTicketServiceClient) {
	fmt.Println("Purchase Ticket Example:")

	req := &train.TicketRequest{
		From:            "London",
		To:              "France",
		UserFirstName:  "John",
		UserLastName:   "Doe",
		UserEmail:      "john.doe@example.com",
	}

	receipt, err := client.PurchaseTicket(context.Background(), req)
	if err != nil {
		log.Fatalf("PurchaseTicket failed: %v", err)
	}

	fmt.Printf("Receipt: %+v\n", receipt)
	fmt.Println()
}

func getReceiptDetailsExample(client train.TrainTicketServiceClient) {
	fmt.Println("Get Receipt Details Example:")

	req := &train.ReceiptRequest{
		UserEmail: "john.doe@example.com",
	}

	receipt, err := client.GetReceiptDetails(context.Background(), req)
	if err != nil {
		log.Fatalf("GetReceiptDetails failed: %v", err)
	}

	fmt.Printf("Receipt Details: %+v\n", receipt)
	fmt.Println()
}

func viewUserSeatsExample(client train.TrainTicketServiceClient) {
	fmt.Println("View User Seats Example:")

	req := &train.UserSeatsRequest{
		SeatSection: "sectionA",
	}

	userSeats, err := client.ViewUserSeats(context.Background(), req)
	if err != nil {
		log.Fatalf("ViewUserSeats failed: %v", err)
	}

	fmt.Printf("User Seats in Section A: %+v\n", userSeats.UserSeats)
	fmt.Println()
}

func removeUserExample(client train.TrainTicketServiceClient) {
	fmt.Println("Remove User Example:")

	req := &train.UserRequest{
		UserEmail: "john.doe@example.com",
	}

	_, err := client.RemoveUser(context.Background(), req)
	if err != nil {
		log.Fatalf("RemoveUser failed: %v", err)
	}

	fmt.Println("User removed successfully.")
	fmt.Println()
}

func modifyUserSeatExample(client train.TrainTicketServiceClient) {
	fmt.Println("Modify User Seat Example:")

	req := &train.ModifyUserSeatRequest{
		UserEmail:       "john.doe@example.com",
		NewSeatSection:  "sectionB",
	}

	receipt, err := client.ModifyUserSeat(context.Background(), req)
	if err != nil {
		log.Fatalf("ModifyUserSeat failed: %v", err)
	}

	fmt.Printf("Modified Receipt: %+v\n", receipt)
	fmt.Println()
}
