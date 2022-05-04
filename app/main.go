package main

import (
	"app/models"
	"fmt"
	"log"
)

func init() {
	log.SetPrefix("Blockchain: ")
}

func main() {
	walletM := models.NewWallet()
	walletA := models.NewWallet()
	walletB := models.NewWallet()

	t := models.NewWalletTransaction(walletA.PrivateKey(), walletA.PublicKey(), walletA.BlockchainAddress(), walletB.BlockchainAddress(), 1.0)

	blockchain := models.NewBlockchain(walletM.BlockchainAddress())
	blockchain.AddTransaction(walletA.PublicKey(), t.GenerateSignature(), walletA.BlockchainAddress(), walletB.BlockchainAddress(), 1.0)

	blockchain.Mining()

	fmt.Printf("A %.1f\n", blockchain.CalculateTotalAmount(walletA.BlockchainAddress()))
	fmt.Printf("B %.1f\n", blockchain.CalculateTotalAmount(walletB.BlockchainAddress()))
	fmt.Printf("M %.1f\n", blockchain.CalculateTotalAmount(walletM.BlockchainAddress()))
}
