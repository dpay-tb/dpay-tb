package main

import (
	"fmt"
	"dpay/transaction"
)



func main() {
	client := transaction.CreateClient()
	defer client.Close()

	client.CreateAccount(transaction.IdFromHex("ab12"))
	fmt.Println("Success.")
}
