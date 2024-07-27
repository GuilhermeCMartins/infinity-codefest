package user

type Currency string

const (
	IC  Currency = "IC"
	USD Currency = "USD"
	BRL Currency = "BRL"
	ETH Currency = "ETH"
	BTC Currency = "BTC"
)

type UserStatus string

const (
	REVIEW   UserStatus = "review"
	SUCCESS  UserStatus = "success"
	FAILED   UserStatus = "failed"
	APPROVED UserStatus = "approved"
)

type UserEvents string

const (
	REQUEST UserEvents = "User.Request"
	PENDING UserEvents = "User.Pending"
	CREATED UserEvents = "User.Created"
)
