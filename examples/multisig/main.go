// Package multisig shows how to create & join multi-signature wallet and send transaction
package main

import (
	"fmt"
	"log"

	bws "github.com/pavel-main/bws-go/client"
	"github.com/pavel-main/bws-go/config"
	"github.com/pavel-main/bws-go/credentials"
	"github.com/pavel-main/bws-go/models"
)

// Actor represents a person with a private key
type Actor struct {
	Name        string
	Client      *bws.Client
	Credentials *credentials.Credentials
}

// NewActor creates new person with specified name
func NewActor(name string, cfg *config.Config) *Actor {
	a := new(Actor)
	a.Name = name
	a.Credentials, _ = credentials.New(cfg, 256)
	a.Client, _ = bws.New(cfg, a.Credentials)
	return a
}

func main() {
	// Init config
	cfg := config.NewPublicTestnet()
	cfg.Debug = false

	// Init actors
	irene := NewActor("Irene", cfg)
	tomas := NewActor("Tomas", cfg)

	// Create wallet
	wallet, err := irene.Client.CreateWallet(irene.Name, 2, 2, false)
	if err != nil {
		log.Fatalf("Error creating wallet: %s", err.Error())
	}

	// Join wallet by Irene
	if _, err := irene.Client.JoinWallet(irene.Name, wallet.Secret); err != nil {
		log.Fatalf("Error joining wallet: %s", err.Error())
	}

	// Join wallet by Tomas
	if _, err := tomas.Client.JoinWallet(tomas.Name, wallet.Secret); err != nil {
		log.Fatalf("Error joining wallet: %s", err.Error())
	}

	// Print credentials
	fmt.Printf("Irene's master key: %s\n", irene.Credentials.RootKey.String())
	fmt.Printf("Tomas's master key: %s\n", tomas.Credentials.RootKey.String())

	// Create tx proposal by Irene
	output := models.NewTxOutputSingle(333333, "mnv9rH2VfAUX9YZzFkoRysGFtggvz1wRnY")
	txp, err := irene.Client.CreateTxProposal(output, "normal", false)
	if err != nil {
		log.Fatalf("Error creating tx proposal: %s", err.Error())
	}

	// Publish tx proposal by Irene
	if _, err := irene.Client.PublishTxProposal(txp); err != nil {
		log.Printf("Error publishing tx proposal: %s", err.Error())
	}

	// Sign tx proposal by Irene
	if _, err := irene.Client.SignTxProposal(txp); err != nil {
		log.Printf("Error signing tx proposal: %s", err.Error())
	}

	// Sign tx proposal by Tomas
	if _, err := tomas.Client.SignTxProposal(txp); err != nil {
		log.Printf("Error signing tx proposal: %s", err.Error())
	}

	// Broadcast tx proposal by Tomas
	txp, err = tomas.Client.BroadcastTxProposal(txp.ID)
	if err != nil {
		log.Printf("Error broadcasting tx proposal: %s", err.Error())
	}

	// Print information
	fmt.Printf("Transaction ID: %s\n", txp.ID)
	fmt.Printf("Transaction status: %s\n", txp.Status)
	fmt.Printf("Transaction hash: %s\n", txp.TxID)
}
