// Package simple shows how to open wallet with a single copayer and send transaction
package main

import (
	"fmt"
	"log"

	bws "github.com/pavel-main/bws-go/client"
	"github.com/pavel-main/bws-go/config"
	"github.com/pavel-main/bws-go/credentials"
	"github.com/pavel-main/bws-go/models"
)

var rootKey = "tprv8ZgxMBicQKsPf1Zu9VstrFcmfHVRBibGLcTKn4ZxEYZkxR8fzUQsj1B49LRze1JpL2GAkL5GbqingWSqcW3cNNngt736xpeLJbYE6mHjaRr"

func main() {
	// Init config
	cfg := config.NewPublicTestnet()
	cfg.Debug = false

	// Init credentials from private key string
	credentials, err := credentials.NewFromPrivateKey(cfg, rootKey)
	if err != nil {
		log.Fatalf("Error creating credentials: %s", err.Error())
	}

	// Init BWS client
	client, err := bws.New(cfg, credentials)
	if err != nil {
		log.Fatalf("Error creating client: %s", err.Error())
	}

	// Create tx proposal
	output := models.NewTxOutputSingle(300000, "mnv9rH2VfAUX9YZzFkoRysGFtggvz1wRnY")
	txp, err := client.CreateTxProposal(output, "normal", false)
	if err != nil {
		log.Fatalf("Error creating tx proposal: %s", err.Error())
	}

	// Publish tx proposal
	if _, err := client.PublishTxProposal(txp); err != nil {
		log.Printf("Error publishing tx proposal: %s", err.Error())
	}

	// Sign tx proposal
	if _, err := client.SignTxProposal(txp); err != nil {
		log.Printf("Error signing tx proposal: %s", err.Error())
	}

	// Broadcast tx proposal
	txp, err = client.BroadcastTxProposal(txp.ID)
	if err != nil {
		log.Printf("Error signing tx proposal: %s", err.Error())
	}

	// Print information
	fmt.Printf("Transaction ID: %s\n", txp.ID)
	fmt.Printf("Transaction status: %s\n", txp.Status)
	fmt.Printf("Transaction hash: %s\n", txp.TxID)
}
