package transaction

var TransactionClient Client = CreateClient()

func Init() {
	go TransferWorker()
}