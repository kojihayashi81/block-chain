package models

import (
	"app/clock"
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

type Block struct {
	nonce        int
	previousHash [32]byte
	timestamp    int64
	transactions []string
}

func NewBlock(nonce int, previousHash [32]byte) *Block {
	b := new(Block)
	b.timestamp = clock.Now().UnixNano()
	b.nonce = nonce
	b.previousHash = previousHash
	b.transactions = []string{}
	return b
}

func (b *Block) Hash() [32]byte {
	m, _ := json.Marshal(b)
	return sha256.Sum256([]byte(m))
}

func (b *Block) Print() {
	fmt.Printf("nonce        :%d\n", b.nonce)
	fmt.Printf("previousHash :%x\n", b.previousHash)
	fmt.Printf("timestamp    :%d\n", b.timestamp)
	fmt.Printf("transactions :%s\n", b.transactions)
}

func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Timestamp    int64    `json:"timestamp"`
		Nonce        int      `json:"nonce"`
		PreviousHash [32]byte `json:"previous_hash"`
		Transactions []string `json:"transaction"`
	}{
		Timestamp:    b.timestamp,
		Nonce:        b.nonce,
		PreviousHash: b.previousHash,
		Transactions: b.transactions,
	})
}
