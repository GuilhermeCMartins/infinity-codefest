package models

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	Id        uuid.UUID          `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Sender    uuid.UUID          `json:"sender" gorm:"type:uuid"`
	Receiver  uuid.UUID          `json:"receiver" gorm:"type:uuid"`
	Amount    float32            `json:"amount" validate:"required,gt=0"`
	Currency  *Currency          `json:"currency"`
	Hash      string             `json:"hash"`
	Status    *TransactionStatus `json:"status"`
	Reason    string             `json:"reason"`
	Users     []*User            `json:"user" gorm:"many2many:transaction_user"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

type TransactionPayload struct {
	Id        uuid.UUID          `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Event     TransactionEvent   `json:"event"`
	Sender    uuid.UUID          `json:"sender" gorm:"type:uuid" validate:"required"`
	Receiver  uuid.UUID          `json:"receiver" gorm:"type:uuid" validate:"required"`
	Amount    float32            `json:"amount" validate:"required,gt=0"`
	Currency  *Currency          `json:"currency"`
	Hash      string             `json:"hash"`
	Status    *TransactionStatus `json:"status"`
	Reason    string             `json:"reason"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}
