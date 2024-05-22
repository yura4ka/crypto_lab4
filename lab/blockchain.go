package lab

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strings"
)

var GENESIS = "<Genesis>"

type Blockchain struct {
	Chain      []*Block
	difficulty int
	Filename   string
}

func NewBlockchain(difficulty int, filename string) Blockchain {
	chain := make([]*Block, 1)
	chain[0] = NewBlock(0, []Transaction{{GENESIS, GENESIS, 0}}, "", 0)
	return Blockchain{chain, difficulty, filename}
}

func FromFile(filename string) (Blockchain, error) {
	b := NewBlockchain(5, filename)
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		return b, err
	}

	err = json.Unmarshal(data, &b)
	return b, err
}

func (b *Blockchain) Save() {
	file, err := os.OpenFile(b.Filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	data, err := json.Marshal(b)
	if err != nil {
		fmt.Println(err)
	}
	file.Write(data)
}

func (b *Blockchain) GetLastBlock() *Block {
	return b.Chain[len(b.Chain)-1]
}

func (b *Blockchain) AddBlock(transactions []Transaction) {
	created := NewBlock(len(b.Chain), transactions, b.GetLastBlock().Hash, 0)
	created.Mine(b.difficulty)
	b.Chain = append(b.Chain, created)
	b.Save()
}

func (b *Blockchain) CheckValid() bool {
	for i := range len(b.Chain) - 1 {
		cur := b.Chain[i+1]
		prev := b.Chain[i]
		if cur.PreviosHash != prev.calculateBlockHash() {
			fmt.Printf("Invalid block %d. cur.prevHash(%s) != prev.hash(%s)\n", i+1, cur.PreviosHash, prev.Hash)
			return false
		}
		if !strings.HasSuffix(cur.Hash, strings.Repeat("0", b.difficulty)) {
			fmt.Printf("Invalid block %d. cur.Hash(%s) has wrong suffix(%d)\n", i+1, cur.Hash, b.difficulty)
			return false
		}
	}
	return true
}

func (bc *Blockchain) GetBalance(user string) int {
	result := 0
	for _, b := range bc.Chain {
		for _, t := range b.Transactions {
			if t.From == user {
				result -= t.Amount
			} else if t.To == user {
				result += t.Amount
			}
		}
	}
	return result
}

func (bc *Blockchain) GetUsers() []string {
	usersMap := make(map[string]bool)
	for _, b := range bc.Chain {
		for _, t := range b.Transactions {
			usersMap[t.From] = true
			usersMap[t.To] = true
		}
	}

	delete(usersMap, GENESIS)
	result := make([]string, 0, len(usersMap))
	for u := range usersMap {
		result = append(result, u)
	}

	return result
}

func (bc *Blockchain) GetUserMinMax(user string) (int, int) {
	min := math.Inf(1)
	max := math.Inf(-1)
	balance := 0
	for _, b := range bc.Chain {
		for _, t := range b.Transactions {
			if user == t.From {
				balance -= t.Amount
			} else if user == t.To {
				balance += t.Amount
			}
			min = math.Min(min, float64(balance))
			max = math.Max(max, float64(balance))
		}
	}

	return int(min), int(max)
}
