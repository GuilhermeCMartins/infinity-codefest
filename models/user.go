package models

import (
	"encoding/json"
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

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

// Balance vira bigint, e criamos um transforamr para guardar em centavos!
type User struct {
	Id        uuid.UUID   `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string      `json:"name"`
	Email     string      `json:"email"`
	PublicKey string      `json:"public_key"`
	Status    *UserStatus `json:"status"`
	Balance   float32     `json:"balance" `
	Currency  *Currency   `json:"currency"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type UserPayload struct {
	Event     UserEvents `json:"event" validate:"required"`
	Name      string     `json:"name" validate:"required"`
	Email     string     `json:"email" validate:"required,email"`
	PublicKey string     `json:"public_key" validate:"required"`
	Balance   float32    `json:"balance" validate:"required"`
	Currency  Currency   `json:"currency" validate:"required"`
	CreatedAt time.Time  `json:"created_at" validate:"required"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// Voltar o User
func CreateUser(db *gorm.DB, user *User) error {
	return db.Create(user).Error
}

func verifyIfCreationIsValid(payload UserPayload) error {
	validate := validator.New()
	err := validate.Struct(payload)
	if err != nil {
		return err
	}
	return nil
}

func handleRequestUser(db *gorm.DB, payload UserPayload) string {
	err := verifyIfCreationIsValid(payload)

	status := REVIEW

	//Ta errado, tem que buscar CreatedAt e UpdatedAt do payload
	user := User{
		Id:        uuid.New(),
		Name:      payload.Name,
		Email:     payload.Email,
		PublicKey: payload.PublicKey,
		Status:    &status,
		Balance:   payload.Balance,
		Currency:  &payload.Currency,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err != nil {
		failedStatus := FAILED
		user.Status = &failedStatus
		log.Printf("Validation failed: %v", err)

		// Mesmo se rejeitar, temos que mandar a mensagem e a reason para nosso broker
		return ""
	}

	err = CreateUser(db, &user)
	if err != nil {
		log.Printf("Failed to create user: %v", err)
		return ""
	}

	message := struct {
		Id        uuid.UUID   `json:"id"`
		Status    *UserStatus `json:"status"`
		Event     UserEvents  `json:"event" validate:"required"`
		Name      string      `json:"name" validate:"required"`
		Email     string      `json:"email" validate:"required,email"`
		PublicKey string      `json:"public_key" validate:"required"`
		Balance   float32     `json:"balance" validate:"required"`
		Currency  Currency    `json:"currency" validate:"required"`
		CreatedAt time.Time   `json:"created_at" validate:"required"`
		UpdatedAt time.Time   `json:"updated_at"`
	}{
		Id:        user.Id,
		Status:    user.Status,
		Event:     PENDING,
		Name:      payload.Name,
		Email:     payload.Email,
		PublicKey: payload.PublicKey,
		Balance:   payload.Balance,
		Currency:  payload.Currency,
		CreatedAt: payload.CreatedAt,
		UpdatedAt: payload.UpdatedAt,
	}

	messageJSON, err := json.Marshal(message)

	if err != nil {
		panic(err)
	}

	var stringMessage = string(messageJSON)

	return stringMessage
}

func HandleMessageUser(db *gorm.DB, payload UserPayload) string {

	var result string

	switch payload.Event {
	case REQUEST:
		result = handleRequestUser(db, payload)
	case PENDING:
	case CREATED:
	default:
		log.Printf("Unknown event type: %s", payload.Event)
	}

	return result
}
