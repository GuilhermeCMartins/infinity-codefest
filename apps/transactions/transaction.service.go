package transaction

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func createMessage(transaction Transaction, event TransactionEvent) string {
	message := struct {
		Id        uuid.UUID          `json:"id"`
		Status    *TransactionStatus `json:"status"`
		Event     TransactionEvent   `json:"event" validate:"required"`
		Amount    float32            `json:"amount" validate:"required"`
		Currency  Currency           `json:"currency" validate:"required"`
		CreatedAt time.Time          `json:"created_at" validate:"required"`
		UpdatedAt time.Time          `json:"updated_at"`
	}{
		Id:        transaction.Id,
		Status:    transaction.Status,
		Event:     event,
		Amount:    transaction.Amount,
		Currency:  *transaction.Currency,
		CreatedAt: transaction.CreatedAt,
		UpdatedAt: transaction.UpdatedAt,
	}

	messageJSON, err := json.Marshal(message)
	if err != nil {
		fmt.Errorf("Failed to marshal message: %v", err)
		return ""
	}

	return string(messageJSON)
}

func verifyIfCreationIsValid(payload TransactionPayload) error {
	validate := validator.New()
	err := validate.Struct(payload)
	if err != nil {
		return err
	}
	return nil
}

// TO-DO: Arrumar os updates e os retornos de mensagens
func handleRequestTransaction(db *gorm.DB, payload TransactionPayload) string {
	err := verifyIfCreationIsValid(payload)

	status := REVIEW

	transaction := Transaction{
		Id:        uuid.New(),
		Status:    &status,
		CreatedAt: payload.CreatedAt,
		UpdatedAt: payload.UpdatedAt,
	}

	if err != nil {
		failedStatus := FAILED

		updates := Transaction{
			Status:    &failedStatus,
			Reason:    "Falta de campos para criação de usuário",
			UpdatedAt: time.Now(),
		}

		transactionUpdated, _ := updateTransaction(db, transaction.Id, updates)

		message := createMessage(transactionUpdated, PENDING)

		return message
	}

	err = createTransaction(db, &transaction)
	if err != nil {
		log.Printf("Failed to create transaction: %v", err)
		return ""
	}

	message := struct {
		Id        uuid.UUID          `json:"id"`
		Status    *TransactionStatus `json:"status"`
		Event     TransactionEvent   `json:"event" validate:"required"`
		CreatedAt time.Time          `json:"created_at" validate:"required"`
		UpdatedAt time.Time          `json:"updated_at"`
	}{
		Id:        transaction.Id,
		Status:    transaction.Status,
		Event:     PENDING,
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

// TO-DO: Arrumar os updates e os retornos de mensagens
func handlePendingTransaction(db *gorm.DB, payload TransactionPayload) string {
	err := verifyIfCreationIsValid(payload)
	if err != nil {
		log.Printf("Validation failed: %v", err)
		return ""
	}

	var transaction Transaction
	result := db.First(&transaction, "id = ?", payload.Id)
	if result.Error != nil {
		log.Printf("Transaction not found: %v", result.Error)
		return ""
	}

	status := APPROVED

	if *payload.Status == FAILED {
		status := FAILED

		updates := Transaction{
			Status:    &status,
			Reason:    "Reprovado pelo KYC/FRAUD",
			UpdatedAt: time.Now(),
		}

		transactionUpdated, _ := updateTransaction(db, transaction.Id, updates)
		//tratar erro de banco

		message := createMessage(transactionUpdated, PENDING)
		return message
	}

	updates := Transaction{
		Status:    &status,
		Reason:    payload.Reason,
		UpdatedAt: time.Now(),
	}

	transactionUpdated, _ := updateTransaction(db, transaction.Id, updates)
	//tratar erro de banco

	message := createMessage(transactionUpdated, CREATED)
	return message
}

func HandleMessageTransaction(db *gorm.DB, payload TransactionPayload) string {

	var result string

	switch payload.Event {
	case REQUEST:
		result = handleRequestTransaction(db, payload)
	case PENDING:
		result = handlePendingTransaction(db, payload)
	default:
		log.Printf("Unknown event type: %s", payload.Event)
	}

	return result
}
