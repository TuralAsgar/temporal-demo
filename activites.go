package main

import (
	"context"
	"log"
)

func PickLeadActivity(ctx context.Context, lead Lead) error {
	log.Println("Picking lead with id:", lead.ID)
	return nil
}

func CallCustomerActivity(ctx context.Context, lead Lead) (string, error) {
	log.Println("Calling customer via 3CX:", lead.Phone)
	// return "", errors.New("3cx is not working")
	return "call-123", nil
}

func UpdateLeadActivity(ctx context.Context, leadID string, result string) error {
	log.Println("Updating lead", leadID, "with result:", result)
	return nil
}

func CreateInteractionActivity(ctx context.Context, leadID string, result string) error {
	log.Println("Creating interaction for lead", leadID)
	// return "", errors.New("Interaction service is not working")
	return nil
}
