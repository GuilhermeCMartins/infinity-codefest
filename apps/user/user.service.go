package user

import (
	"encoding/json"
	"fmt"
	"log"
	apps "myapp/apps/transactions"
	"time"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TO-DO:
func createMessage(user User, event UserEvents) string {
	message := struct {
		Id        uuid.UUID   `json:"id"`
		Status    *UserStatus `json:"status"`
		Event     UserEvents  `json:"event" validate:"required"`
		Name      string      `json:"name" validate:"required"`
		Email     string      `json:"email" validate:"required,email"`
		PublicKey string      `json:"public_key" validate:"required"`
		Balance   float64     `json:"balance" validate:"required"`
		Currency  Currency    `json:"currency" validate:"required"`
		CreatedAt time.Time   `json:"created_at" validate:"required"`
		UpdatedAt time.Time   `json:"updated_at"`
	}{
		Id:        user.Id,
		Status:    user.Status,
		Event:     event,
		Name:      user.Name,
		Email:     user.Email,
		PublicKey: user.PublicKey,
		Balance:   user.Balance,
		Currency:  *user.Currency,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	messageJSON, err := json.Marshal(message)
	if err != nil {
		fmt.Errorf("Failed to marshal message: %v", err)
		return ""
	}

	return string(messageJSON)
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

	user := User{
		Id:        uuid.New(),
		Name:      payload.Name,
		Email:     payload.Email,
		PublicKey: payload.PublicKey,
		Status:    &status,
		Balance:   payload.Balance,
		Currency:  &payload.Currency,
		CreatedAt: payload.CreatedAt,
		UpdatedAt: payload.UpdatedAt,
	}

	if err != nil {
		failedStatus := FAILED

		updates := User{
			Status:    &failedStatus,
			Reason:    "Falta de campos para criação de usuário",
			UpdatedAt: time.Now(),
		}

		userUpdated, _ := UpdateUser(db, user.Id, updates)

		message := createMessage(userUpdated, PENDING)

		return message
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
		Balance   float64     `json:"balance" validate:"required"`
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

// TO-DO: verify if message already consumed
func handlePendingUser(db *gorm.DB, payload UserPayload) string {
	err := verifyIfCreationIsValid(payload)
	if err != nil {
		log.Printf("Validation failed: %v", err)
		return ""
	}

	var user User
	result := db.First(&user, "id = ?", payload.Id)
	if result.Error != nil {
		log.Printf("User not found: %v", result.Error)
		return ""
	}

	status := APPROVED

	if *payload.Status == FAILED {
		status := FAILED

		updates := User{
			Status:    &status,
			Reason:    "Reprovado pelo KYC/FRAUD",
			UpdatedAt: time.Now(),
		}

		userUpdated, _ := UpdateUser(db, user.Id, updates)
		//tratar erro de banco

		message := createMessage(userUpdated, PENDING)
		return message
	}

	updates := User{
		Status:    &status,
		Reason:    payload.Reason,
		UpdatedAt: time.Now(),
	}

	userUpdated, _ := UpdateUser(db, user.Id, updates)
	//tratar erro de banco

	message := createMessage(userUpdated, CREATED)
	println("[USER PENDING]", message)
	return message
}

func HandleMessageUser(db *gorm.DB, payload UserPayload) string {

	var result string

	switch payload.Event {
	case REQUEST:
		result = handleRequestUser(db, payload)
	case PENDING:
		result = handlePendingUser(db, payload)
	default:
		log.Printf("Unknown event type: %s", payload.Event)
	}

	return result
}

type UsersResponse struct {
	Users []User `json:"users"`
	Count int    `json:"count"`
}

func FindAllUsers(db *gorm.DB) (users []User, count int, err error) {
	result := db.Find(&users)
	
	if result.Error != nil {
			return nil, 0, result.Error
	}
	
	count = len(users)
	
	return users, count, nil
}

func FindUserById(db *gorm.DB, id uuid.UUID) (User, error) {
	var user User
	result := db.First(&user, "id = ?", id)
	
	if result.Error != nil {
		return User{}, result.Error
	}
	
	return user, nil
}

func FindUserTransactions(db *gorm.DB, userID uuid.UUID) ([]apps.Transaction, int, error) {
	var transactions []apps.Transaction

	result := db.Where("sender = ? OR receiver = ?", userID, userID).Find(&transactions)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	count := len(transactions)

	return transactions, count, nil
}

// findTransactionsByUserId
// findUserTransactionByStatus /users/:id/transactions/:tx
// findUserTransactionByStatus /users/:id/transactions/status/:status