package models

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
	TX_REVIEW   TransactionStatus = "review"
	TX_SUCCESS  TransactionStatus = "success"
	TX_FAILED   TransactionStatus = "failed"
	TX_APPROVED TransactionStatus = "approved"
)

type TransactionEvent string

const (
	TX_REQUEST TransactionEvent = "Transaction.Request"
	TX_PENDING TransactionEvent = "Transaction.Pending"
	TX_CREATED TransactionEvent = "Transaction.Created"
)

type UserStatus string

const (
	USER_REVIEW   UserStatus = "review"
	USER_SUCCESS  UserStatus = "success"
	USER_FAILED   UserStatus = "failed"
	USER_APPROVED UserStatus = "approved"
)

type UserEvents string

const (
	USER_REQUEST UserEvents = "User.Request"
	USER_PENDING UserEvents = "User.Pending"
	USER_CREATED UserEvents = "User.Created"
)
