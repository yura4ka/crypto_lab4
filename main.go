package main

import (
	"fmt"

	"github.com/yura4ka/crypto_lab4/lab"
)

func printData(bc lab.Blockchain) {
	fmt.Println("Is valid: ", bc.CheckValid())
	users := bc.GetUsers()
	for _, u := range users {
		balance := bc.GetBalance(u)
		min, max := bc.GetUserMinMax(u)
		fmt.Printf("User %s: Balance: %d; min balance: %d, max balance: %d\n", u, balance, min, max)
	}
}

func main() {
	bc := lab.NewBlockchain(5, "bc.json")
	bc.AddBlock([]lab.Transaction{
		{From: "Alice", To: "Bob", Amount: 10},
		{From: "Alice", To: "Bob", Amount: 15},
		{From: "Bob", To: "Mark", Amount: 12},
	})
	bc.AddBlock([]lab.Transaction{
		{From: "Mark", To: "Alice", Amount: 3},
		{From: "Greg", To: "Alice", Amount: 10},
	})
	bc.AddBlock([]lab.Transaction{
		{From: "Greg", To: "Mark", Amount: 12},
		{From: "Bob", To: "Greg", Amount: 8},
	})

	printData(bc)

	fmt.Print("\nLoading from file...\n\n")
	bc2, err := lab.FromFile("bc.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	bc2.AddBlock([]lab.Transaction{{From: "Bob", To: "Alice", Amount: 22}})
	printData(bc2)
}
