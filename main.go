package main

import (
	"fmt"
	"dpay/transaction"
)



func main() {
	client := transaction.CreateClient()
	client.InitializeBank()
	defer client.Close()

	//fmt.Printf("Bank id = %v\n", transaction.BANK_ID)
	id := transaction.IdFromHex("abcdef")
	id2 := transaction.IdFromHex("fedcba")
	//fmt.Printf("Creating account %v...", id)
	client.CreateWithBalance(id, 4000)
	client.CreateWithBalance(id2, 4000)

	client.Transfer(id, id2, 2000)

	fmt.Printf("Balance of %v = %v\n", id, client.GetBalance(id))
	fmt.Printf("Balance of %v = %v\n", id2, client.GetBalance(id2))
	fmt.Println("Success.")
}
