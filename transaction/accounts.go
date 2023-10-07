package transaction

import (
	"os"
	"log"
	"encoding/json"
	tb "github.com/tigerbeetle/tigerbeetle-go"
	tb_types "github.com/tigerbeetle/tigerbeetle-go/pkg/types"
	"github.com/google/uuid"
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

	// Batch size for executing queries to the DB
	BATCH_SIZE = 1024

	// If a batch hasn't been processed in the last
	// BATCH_TICK_INTERVAL ms, process the current
	// batch anyway even if it is not BATCH_SIZE
	BATCH_TICK_INTERVAL = 30

)
// When creating accounts, we use the `debits_must_not_exceed_credits`
// flag. This means that the account can't pay more than its balance.
// But to deposit an initial amount into the account, we need to
// transfer it from a special account, which we call the bank account.
// This will be created with the `credits_must_not_exceed_debits` flag,
// and will only be used to transfer initial amounts to new accounts.
// We're giving a special ID to the bank account.
var BANK_ID = IdFromHex("ffffffffffffffffffffffffffffff00")

type Client struct {
	client tb.Client
}

type AccountId struct {
	id tb_types.Uint128
}

func (a *AccountId) UnmarshalJSON(data []byte) error {
    // Unmarshal the JSON data into a string
    var id string
    if err := json.Unmarshal(data, &id); err != nil {
        return err
    }

    // Convert the string to AccountId using IdFromHex
    *a = IdFromHex(id)

    return nil
}

func IdFromHex(id string) AccountId {
	x, err := tb_types.HexStringToUint128(id)
	if err != nil {
		panic(err)
	}
	return AccountId{x}
}

func NewId() AccountId {
	return AccountId{tb_types.Uint128(uuid.New())}
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
			Ledger:     LEDGER,
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

func (c *Client) InitializeBank() {
	accountsRes, err := c.client.CreateAccounts([]tb_types.Account{
		{
			ID:         BANK_ID.id,
			UserData:   tb_types.Uint128{},
			Reserved:       [48]uint8{},
			Ledger:     LEDGER,
			Code:       DEFAULT_CODE,

			// Bank is used to only transfer initial amounts
			// so its debits always exceeds credits. Ideally
			// the bank should not have any credits
			Flags:      tb_types.AccountFlags{CreditsMustNotExceedDebits: true}.ToUint16(),

			DebitsPending:  0,
			DebitsPosted:   0,
			CreditsPending: 0,
			CreditsPosted:  0,
			Timestamp:  0,
		},
	})
	if err != nil {
		log.Printf("InitializeBank: Error creating accounts: %s", err)
		return
	}
	
	for _, err := range accountsRes {
		log.Printf("IntializeBank: Error creating account %d: %s", err.Index, err.Result)
		return
	}

	log.Println("Bank account created. ")
}

func (c *Client) Transfer(from AccountId, to AccountId, amount uint64) {
	transfer := tb_types.Transfer{
		ID:         tb_types.Uint128(uuid.New()),
		PendingID:      tb_types.Uint128{},
		DebitAccountID:     from.id,
		CreditAccountID:    to.id,
		UserData:       tb_types.Uint128{},
		Reserved:       tb_types.Uint128{},
		Timeout:        0,
		Ledger:         LEDGER,
		Code:           DEFAULT_CODE,
		Flags:          0,
		Amount:         amount,
		Timestamp:      0,
	}
	
	transfersRes, err := c.client.CreateTransfers([]tb_types.Transfer{transfer})
	if err != nil {
		log.Printf("Transfer: Error creating transfer batch: %s", err)
		return
	}

	for _, err := range transfersRes {
		log.Printf("Transfer: Error creating transfer %d: %s", err.Index, err.Result)
		return
	}

	log.Println("Transfer: Transfer successful.")
}

func (c *Client) CreateWithBalance(id AccountId, amount uint64) {
	account := tb_types.Account{
			ID:         id.id,
			UserData:   tb_types.Uint128{},
			Reserved:       [48]uint8{},
			Ledger:     LEDGER,
			Code:       DEFAULT_CODE,


			Flags:      tb_types.AccountFlags{DebitsMustNotExceedCredits: true}.ToUint16(),

			DebitsPending:  0,
			DebitsPosted:   0,
			CreditsPending: 0,
			CreditsPosted:  0,
			Timestamp:  0,
		}

	accountsRes, err := c.client.CreateAccounts([]tb_types.Account{
		account,
	})
	if err != nil {
		log.Printf("CreateWithBalance: Error creating accounts: %s", err)
		return
	}
	
	for _, err := range accountsRes {
		log.Printf("CreateWithBalance: Error creating account %d: %s", err.Index, err.Result)
		return 
	}

	c.Transfer(BANK_ID, id, uint64(amount))
}

func (c *Client) GetAccount(id AccountId) tb_types.Account {
	accounts, err := c.client.LookupAccounts([]tb_types.Uint128{id.id})
	if err != nil {
		panic(err)
	}

	return accounts[0]
}

func (c *Client) GetBalance(id AccountId) uint64 {
	account := c.GetAccount(id)
	return account.CreditsPosted - account.DebitsPosted
}
func (c *Client) Close() {
	c.client.Close()
}