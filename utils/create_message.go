package utils

import (
	"encoding/json"
	"fmt"
	"myapp/models"
	"time"

	"github.com/google/uuid"
)


func CreateMessage(model interface{}, event interface{}) string {
	switch m := model.(type) {
	case models.Transaction:
			message := struct {
					Id        uuid.UUID                 `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
					Sender    uuid.UUID                 `json:"sender" gorm:"type:uuid"`
					Receiver  uuid.UUID                 `json:"receiver" gorm:"type:uuid"`
					Amount    float32                   `json:"amount"`
					Currency  models.Currency           `json:"currency"`
					Hash      string                    `json:"hash"`
					Status    *models.TransactionStatus `json:"status"`
					Reason    string                    `json:"reason"`
					CreatedAt time.Time                 `json:"created_at"`
					UpdatedAt time.Time                 `json:"updated_at"`
					Event     models.TransactionEvent   `json:"event" validate:"required"`
			}{
					Id:        m.Id,
					Status:    m.Status,
					Sender:    m.Sender,
					Receiver:  m.Receiver,
					Hash:      m.Hash,
					Reason:    m.Reason,
					Event:     event.(models.TransactionEvent),
					Amount:    m.Amount,
					Currency:  *m.Currency,
					CreatedAt: m.CreatedAt,
					UpdatedAt: m.UpdatedAt,
			}

			messageJSON, err := json.Marshal(message)
			if err != nil {
					fmt.Errorf("Failed to marshal message: %v", err)
					return ""
			}

			return string(messageJSON)
	case models.User:{
		message := struct {
			Id        uuid.UUID          `json:"id"`
			Status    *models.UserStatus `json:"status"`
			Event     models.UserEvents  `json:"event" validate:"required"`
			Name      string             `json:"name" validate:"required"`
			Email     string             `json:"email" validate:"required,email"`
			PublicKey string             `json:"public_key" validate:"required"`
			Balance   float64            `json:"balance" validate:"required"`
			Currency  models.Currency    `json:"currency" validate:"required"`
			CreatedAt time.Time          `json:"created_at" validate:"required"`
			UpdatedAt time.Time          `json:"updated_at"`
		}{
			Id:        m.Id,
			Status:    m.Status,
			Event:     event.(models.UserEvents),
			Name:      m.Name,
			Email:     m.Email,
			PublicKey: m.PublicKey,
			Balance:   m.Balance,
			Currency:  *m.Currency,
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
		}
	
		messageJSON, err := json.Marshal(message)
		if err != nil {
			fmt.Errorf("Failed to marshal message: %v", err)
			return ""
		}
	
		return string(messageJSON)
	}
	default:
			return ""
	}
}