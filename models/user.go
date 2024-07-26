package models

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Define the Currency type
type Currency string

// Define the possible values for Currency
const (
  IC  Currency = "IC"
  USD Currency = "USD"
  BRL Currency = "BRL"
  ETH Currency = "ETH"
  BTC Currency = "BTC"
)

// Define the UserStatus type
type UserStatus string

// Define the possible values for UserStatus
const (
  Review  UserStatus = "review"
  Success UserStatus = "success"
  Failed  UserStatus = "failed"
  Approved UserStatus = "approved"
)

type User struct {
  Id uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
  Name  string `json:"name"`
  Email  string `json:"email"` 
  PublicKey string `json:"public_key"`
  Status   *UserStatus `json:"status"`
  Balance   decimal.Decimal `json:"balance" gorm:"type:decimal(20,2)"`
  Currency *Currency    `json:"currency"`
  CreatedAt int    `json:"created_at"`
  UpdatedAt int    `json:"updated_at"`
}