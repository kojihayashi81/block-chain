package app

import (
	"app/clock"
	"fmt"
	"log"
)

type App interface {
	NewBlock(nonce int, previousHash string) *Block
	NewBlockChain() *Blockchain
	CreateBlock(nonce int, previousHash string) *Block
	Print()
}

type Block struct {
	nonce        int
	previousHash string
	timestamp    int64
	transactions []string
}

type Blockchain struct {
	transactionPool []string
	chain           []*Block
}

func NewBlock(nonce int, previousHash string) *Block {
	b := new(Block)
	b.timestamp = clock.Now().UnixNano()
	b.nonce = nonce
	b.previousHash = previousHash
	b.transactions = []string{}
	return b
}

func NewBlockChain() *Blockchain {
	bc := new(Blockchain)
	bc.CreateBlock(0, "init hash")
	return bc
}

func (bc *Blockchain) CreateBlock(nonce int, previousHash string) *Block {
	b := NewBlock(nonce, previousHash)
	bc.chain = append(bc.chain, b)
	return b
}

func (b *Block) Print() {
	fmt.Printf("nonce        :%d\n", b.nonce)
	fmt.Printf("previousHash :%s\n", b.previousHash)
	fmt.Printf("timestamp    :%d\n", b.timestamp)
	fmt.Printf("transactions :%s\n", b.transactions)
}

func init() {
	log.SetPrefix("Blockchain: ")
}

func main() {
	b := NewBlock(0, "init hash")
	b.Print()
}
