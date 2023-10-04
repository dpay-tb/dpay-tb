package handlers

import (
	"net/http"
	"dpay/transaction"
	"github.com/gin-gonic/gin"
)

// GetBalance is a handler for the /balance/:id endpoint.
// It returns the balance of the account with the given ID.
func GetBalance(c *gin.Context) {
	id := c.Param("id")
	accountId := transaction.IdFromHex(id)
	balance := transaction.TransactionClient.GetBalance(accountId)

	c.JSON(http.StatusOK, gin.H{"balance": balance})
}

// Transfer is a handler for the /transfer endpoint.
// It transfers the given amount from the source account
// to the destination account.
func Transfer(c *gin.Context) {
	var json struct {
		SourceId string `json:"source_id"`
		DestId string `json:"dest_id"`
		Amount uint64 `json:"amount"`
	}
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sourceId := transaction.IdFromHex(json.SourceId)
	destId := transaction.IdFromHex(json.DestId)
	amount := json.Amount

	transaction.TransactionClient.Transfer(sourceId, destId, amount)

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// CreateWithBalance is a handler for the /create endpoint.
// It creates an account with the given ID and the given balance.
func CreateWithBalance(c *gin.Context) {
	var json struct {
		Id string `json:"id"`
		Amount uint64 `json:"amount"`
	}
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := transaction.IdFromHex(json.Id)
	amount := json.Amount

	transaction.TransactionClient.CreateWithBalance(id, amount)

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}