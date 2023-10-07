package transaction

import (
	//tb "github.com/tigerbeetle/tigerbeetle-go"
	"errors"
	"log"
	"time"
	tb_types "github.com/tigerbeetle/tigerbeetle-go/pkg/types"
	"github.com/google/uuid"
)

var TransferRequests chan TransferQuery = make(chan TransferQuery, BATCH_SIZE)

type TransferRequest struct {
	SourceId AccountId `json:"source_id"`
	DestId AccountId `json:"dest_id"`
	Amount uint64 `json:"amount"`
}

type TransferResponse struct {
	Err error
}

type TransferQuery struct {
	Request TransferRequest
	Response chan TransferResponse
}

func MakeTransferObject(from AccountId, to AccountId, amount uint64) tb_types.Transfer {
	return tb_types.Transfer{
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
}

func DispatchTransers(batch []TransferQuery) {
	transfers := make([]tb_types.Transfer, 0, len(batch))

	for _, query := range batch {
		transfers = append(transfers, MakeTransferObject(query.Request.SourceId, query.Request.DestId, query.Request.Amount))
	}

	transfersRes, err := TransactionClient.client.CreateTransfers(transfers)
	if err != nil {
		log.Fatalf("Transfer: Error creating transfer batch: %s", err)

		// Send the error to all the responses
		for _, query := range batch {
			query.Response <- TransferResponse{Err: err}
		}

		return
	}

	// Make a map of index to errors
	transfer_errors := make(map[int]error)
	for _, err := range transfersRes {
		transfer_errors[int(err.Index)] = errors.New(err.Result.String())
	}

	for i, query := range batch {
		if err, ok := transfer_errors[i]; ok {
			query.Response <- TransferResponse{Err: err}
		} else {
			query.Response <- TransferResponse{Err: nil}
		}
	}

}

func TransferWorker() {
	var batch []TransferQuery
    var batchSize int

    ticker := time.NewTicker(BATCH_TICK_INTERVAL * time.Millisecond)
    for {
        select {
        case req := <-TransferRequests:
            batch = append(batch, req)
            batchSize++
            if batchSize >= BATCH_SIZE {
                DispatchTransers(batch)
                batchSize = 0
                batch = nil
            }
			ticker.Reset(BATCH_TICK_INTERVAL * time.Millisecond)
			
        case <-ticker.C:
            if batchSize > 0 {
                DispatchTransers(batch)
                batchSize = 0
                batch = nil
            }
        }
    }
}