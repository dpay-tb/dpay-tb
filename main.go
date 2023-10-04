package main

import (
	"dpay/transaction"
	"dpay/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	client := transaction.TransactionClient
	defer client.Close()

	router := gin.Default()
	router.GET("/balance/:id", handlers.GetBalance)
	router.POST("/transfer", handlers.Transfer)
	router.POST("/create", handlers.CreateWithBalance)

	router.Run("localhost:8080")
}
