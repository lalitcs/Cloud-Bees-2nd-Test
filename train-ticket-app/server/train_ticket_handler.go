package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"

	train "path-to-proto/train_ticket"
)

var (
	receipts     = make(map[string]*train.Receipt)
	sectionSeats = map[string][]string{
		"sectionA": {"A1", "A2", "A3"},
		"sectionB": {"B1", "B2", "B3"},
	}
	mu sync.Mutex
)

func (s *trainTicketService) PurchaseTicket(ctx context.Context, req *train.TicketRequest) (*train.Receipt, error) {
	mu.Lock()
	defer mu.Unlock()

	// Check if user already exists
	if _, exists := receipts[req.UserEmail]; exists {
		return nil, fmt.Errorf("user with email %s already has a ticket", req.UserEmail)
	}

	// Allocate a seat randomly
	section := getRandomSection()
	seat, err := getAvailableSeat(section)
	if err != nil {
		return nil, fmt.Errorf("failed to allocate seat: %v", err)
	}

	// Create receipt
	receipt := &train.Receipt{
		From:          req.From,
		To:            req.To,
		UserFirstName: req.UserFirstName,
		UserLastName:  req.UserLastName,
		UserEmail:     req.UserEmail,
		PricePaid:     20.0, // Fixed price for the ticket
		SeatSection:   section + seat,
	}

	// Store receipt
	receipts[req.UserEmail] = receipt

	return receipt, nil
}

func (s *trainTicketService) GetReceiptDetails(ctx context.Context, req *train.ReceiptRequest) (*train.Receipt, error) {
	mu.Lock()
	defer mu.Unlock()

	receipt, exists := receipts[req.UserEmail]
	if !exists {
		return nil, fmt.Errorf("user with email %s does not have a ticket", req.UserEmail)
	}

	return receipt, nil
}

func (s *trainTicketService) ViewUserSeats(ctx context.Context, req *train.UserSeatsRequest) (*train.UserSeatsResponse, error) {
	mu.Lock()
	defer mu.Unlock()

	var userSeats []*train.Receipt
	for _, receipt := range receipts {
		if receipt.SeatSection[:1] == req.SeatSection {
			userSeats = append(userSeats, receipt)
		}
	}

	return &train.UserSeatsResponse{UserSeats: userSeats}, nil
}

func (s *trainTicketService) RemoveUser(ctx context.Context, req *train.UserRequest) (*train.EmptyResponse, error) {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := receipts[req.UserEmail]; !exists {
		return nil, fmt.Errorf("user with email %s does not have a ticket", req.UserEmail)
	}

	delete(receipts, req.UserEmail)
	return &train.EmptyResponse{}, nil
}

func (s *trainTicketService) ModifyUserSeat(ctx context.Context, req *train.ModifyUserSeatRequest) (*train.Receipt, error) {
	mu.Lock()
	defer mu.Unlock()

	receipt, exists := receipts[req.UserEmail]
	if !exists {
		return nil, fmt.Errorf("user with email %s does not have a ticket", req.UserEmail)
	}

	// Check if the new seat section is valid
	if _, valid := sectionSeats[req.NewSeatSection]; !valid {
		return nil, fmt.Errorf("invalid seat section: %s", req.NewSeatSection)
	}

	// Allocate a new seat in the requested section
	newSeat, err := getAvailableSeat(req.NewSeatSection)
	if err != nil {
		return nil, fmt.Errorf("failed to allocate new seat: %v", err)
	}

	// Update the receipt with the new seat section and seat
	receipt.SeatSection = req.NewSeatSection + newSeat

	return receipt, nil
}

func getRandomSection() string {
	sections := []string{"sectionA", "sectionB"}
	return sections[rand.Intn(len(sections))]
}

func getAvailableSeat(section string) (string, error) {
	mu.Lock()
	defer mu.Unlock()

	seats, exists := sectionSeats[section]
	if !exists {
		return "", fmt.Errorf("invalid seat section: %s", section)
	}

	for _, seat := range seats {
		if _, occupied := findOccupiedSeat(seat); !occupied {
			return seat, nil
		}
	}

	return "", fmt.Errorf("no available seats in section: %s", section)
}

func findOccupiedSeat(seat string) (*train.Receipt, bool) {
	for _, receipt := range receipts {
		if receipt.SeatSection == seat {
			return receipt, true
		}
	}
	return nil, false
}
