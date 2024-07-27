package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id        uuid.UUID   `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string      `json:"name"`
	Email     string      `json:"email"`
	PublicKey string      `json:"public_key"`
	Status    *UserStatus `json:"status"`
	Balance   float32     `json:"balance" `
	Currency  *Currency   `json:"currency"`
	Reason    string      `json:"reason"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type UserPayload struct {
	Id        uuid.UUID   `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Event     UserEvents  `json:"event" validate:"required"`
	Name      string      `json:"name" validate:"required"`
	Email     string      `json:"email" validate:"required,email"`
	PublicKey string      `json:"public_key" validate:"required"`
	Balance   float32     `json:"balance" validate:"required"`
	Currency  Currency    `json:"currency" validate:"required"`
	Status    *UserStatus `json:"status"`
	Reason    string      `json:"reason"`
	CreatedAt time.Time   `json:"created_at" validate:"required"`
	UpdatedAt time.Time   `json:"updated_at"`
}