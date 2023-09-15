package transaction

import (
	"os"
	"log"
	tb "github.com/tigerbeetle/tigerbeetle-go"
	tb_types "github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

const (
	// This is the default ledger used in the
	// tigerbeetle docs, so we just use the same
	LEDGER = 700

	// Since we're running a single node cluster
	// We can just hardcode it for now
	// In case we decide to use a multi-node cluster
	// this can be changed.
	CLUSTER_ID = 0

	MAX_CONCURRENCY = 32

	// Code is used to differentiate between different
	// types of accounts. For now, we just use the default
	// given in the docs.
	DEFAULT_CODE = 718
)

type Client struct {
	client tb.Client
}

type AccountId struct {
	id tb_types.Uint128
}

func IdFromHex(id string) AccountId {
	x, err := tb_types.HexStringToUint128(id)
	if err != nil {
		panic(err)
	}
	return AccountId{x}
}

func CreateClient() Client {
	port := os.Getenv("TB_ADDRESS")
	if port == "" {
		port = "3000"
	}
	client, err := tb.NewClient(CLUSTER_ID, []string{port}, MAX_CONCURRENCY)

	if err != nil {
		panic(err)
	}

	log.Printf("Client created: %v", client)
	return Client{client}
}

func (c *Client) CreateAccount(id AccountId) {
	accountsRes, err := c.client.CreateAccounts([]tb_types.Account{
		{
			ID:         id.id,
			UserData:   tb_types.Uint128{},
			Reserved:       [48]uint8{},
			Ledger:     1,
			Code:       DEFAULT_CODE,

			// This flag must be set so that the account
			// can't pay more than its balance
			Flags:      tb_types.AccountFlags{DebitsMustNotExceedCredits: true}.ToUint16(),

			DebitsPending:  0,
			DebitsPosted:   0,
			CreditsPending: 0,
			CreditsPosted:  0,
			Timestamp:  0,
		},
	})
	if err != nil {
		log.Printf("Error creating accounts: %s", err)
		return
	}
	
	for _, err := range accountsRes {
		log.Printf("Error creating account %d: %s", err.Index, err.Result)
		return
	}
}

func (c *Client) Close() {
	c.client.Close()
}