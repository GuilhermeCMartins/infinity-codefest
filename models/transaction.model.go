package models

import (
	"time"

	"github.com/google/uuid"
)

	type Transaction struct {
		Id        uuid.UUID          `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
		Sender    string             `json:"sender"`
		Receiver  string             `json:"receiver"`
		Amount    float32            `json:"amount"`
		Currency  *Currency          `json:"currency"`
		Hash      string             `json:"hash"`
		Status    *TransactionStatus `json:"status"`
		Reason    string             `json:"reason"`
		CreatedAt time.Time          `json:"created_at"`
		UpdatedAt time.Time          `json:"updated_at"`
	}

	type TransactionPayload struct {
		Id        uuid.UUID          `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
		Event     TransactionEvent   `json:"event"`
		Sender    string             `json:"sender" validate:"required"`
		Receiver  string             `json:"receiver" validate:"required"`
		Amount    float32            `json:"amount" validate:"required,gt=0"`
		Currency  *Currency          `json:"currency"`
		Hash      string             `json:"hash"`
		Status    *TransactionStatus `json:"status"`
		Reason    string             `json:"reason"`
		CreatedAt time.Time          `json:"created_at"`
		UpdatedAt time.Time          `json:"updated_at"`
	}
