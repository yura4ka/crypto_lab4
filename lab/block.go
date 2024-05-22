package lab

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Transaction struct {
	From   string
	To     string
	Amount int
}

type Block struct {
	Index        int
	Timestamp    time.Time
	Transactions []Transaction
	PreviosHash  string
	Nonce        int
	MerkleRoot   TreeNode
	Hash         string
}

func (b *Block) calculateBlockHash() string {
	b.Hash = ""
	v, _ := json.Marshal(b)
	hash := sha256.New()
	hash.Write(v)
	b.Hash = hex.EncodeToString(hash.Sum(nil))
	return b.Hash
}

func NewBlock(index int, transactions []Transaction, previousHash string, nonce int) *Block {
	b := Block{index, time.Now(), transactions, previousHash, nonce, BuildMerkleTree(transactions), ""}
	b.calculateBlockHash()
	return &b
}

func (b *Block) Mine(difficulty int) {
	fmt.Println("mining block", b.Index)
	target := strings.Repeat("0", difficulty)
	for !strings.HasSuffix(b.Hash, target) {
		b.Nonce += 1
		b.calculateBlockHash()
	}
}
