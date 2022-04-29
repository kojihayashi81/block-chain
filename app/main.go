package main

import (
	"app/models"
	"log"
)

// type App interface {
// 	NewBlock(nonce int, previousHash string) *models.Block
// 	NewBlockChain() *models.Blockchain
// 	CreateBlock(nonce int, previousHash string) *models.Block
// 	Print()
// }

func init() {
	log.SetPrefix("Blockchain: ")
}

func main() {
	blockChain := models.NewBlockChain()
	blockChain.Print()
}
