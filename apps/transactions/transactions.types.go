package transactions

type Currency string

const (
	IC  Currency = "IC"
	USD Currency = "USD"
	BRL Currency = "BRL"
	ETH Currency = "ETH"
	BTC Currency = "BTC"
)

type TransactionStatus string

const (
	REVIEW   TransactionStatus = "review"
	SUCCESS  TransactionStatus = "success"
	FAILED   TransactionStatus = "failed"
	APPROVED TransactionStatus = "approved"
)

type TransactionEvent string

const (
	REQUEST TransactionEvent = "Transaction.Request"
	PENDING TransactionEvent = "Transaction.Pending"
	CREATED TransactionEvent = "Transaction.Created"
)
