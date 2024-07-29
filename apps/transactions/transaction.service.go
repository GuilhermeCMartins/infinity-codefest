package transaction

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"myapp/apps/user"
	"time"

	"myapp/utils"

	eth "github.com/ethereum/go-ethereum/crypto"
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

func validateHash(payload TransactionPayload, publicKey string) (bool, error) {
	data := fmt.Sprintf("%s%s%f%s%s", payload.Sender, payload.Receiver, payload.Amount, payload.CreatedAt.String(), payload.Currency)

	hash := eth.Keccak256([]byte(data))

	pubKeyBytes, err := hex.DecodeString(publicKey)
	if err != nil {
		return false, fmt.Errorf("failed to decode public key: %v", err)
	}
	pubKey, err := eth.UnmarshalPubkey(pubKeyBytes)
	if err != nil {
		return false, fmt.Errorf("failed to unmarshal public key: %v", err)
	}

	sigBytes, err := hex.DecodeString(payload.Hash)
	if err != nil {
		return false, fmt.Errorf("failed to decode signature: %v", err)
	}

	r := big.Int{}
	s := big.Int{}
	sigLen := len(sigBytes)
	r.SetBytes(sigBytes[:(sigLen / 2)])
	s.SetBytes(sigBytes[(sigLen / 2):])

	isValid := ecdsa.Verify(pubKey, hash, &r, &s)

	return isValid, nil
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
	error := verifyIfCreationIsValid(payload)

	status := REVIEW

	var requested user.User
	if err := db.First(&requested, "public_key = ?", payload.Sender).Error; err != nil {
		log.Printf("User not found: %v", err)
		return ""
	}

	var receiver user.User
	if err := db.First(&requested, "public_key = ?", payload.Receiver).Error; err != nil {
		log.Printf("User not found: %v", err)
		return ""
	}

	valid, errHash := validateHash(payload, requested.PublicKey)
	if errHash != nil {
		log.Printf("Hash validation failed: %v", errHash)
		return ""
	}

	if !valid {
		log.Printf("Invalid hash")
		failedStatus := FAILED

		updates := Transaction{
			Status:    &failedStatus,
			Reason:    "Invalid hash",
			UpdatedAt: time.Now(),
		}

		updateTransaction(db, payload.Id, updates)

		return ""
	}

	senderAmount, _, errConvert := utils.ConvertCurrency(float64(payload.Amount), utils.Currency(*requested.Currency), utils.Currency(*receiver.Currency), utils.Currency(*payload.Currency))
	if errConvert != nil {
		log.Printf("Conversion failed: %v", errConvert)
		return ""
	}

	hasValue := senderAmount - float64(requested.Balance)

	if hasValue < 0 {
		log.Printf("Insuficient balance")
		return ""
	}

	transaction := Transaction{
		Id:        uuid.New(),
		Status:    &status,
		Sender:    payload.Sender,
		Receiver:  payload.Receiver,
		Amount:    payload.Amount,
		Currency:  payload.Currency,
		Hash:      payload.Hash,
		CreatedAt: payload.CreatedAt,
		UpdatedAt: payload.UpdatedAt,
	}

	if error != nil {
		failedStatus := FAILED

		updates := Transaction{
			Status:    &failedStatus,
			Reason:    "Falta de campos para criação de transação",
			UpdatedAt: time.Now(),
		}

		transactionUpdated, _ := updateTransaction(db, transaction.Id, updates)

		message := createMessage(transactionUpdated, PENDING)

		return message
	}

	if err := createTransaction(db, &transaction); err != nil {
		log.Printf("Failed to create transaction: %v", err)
		return ""
	}

	message := struct {
		Id        uuid.UUID          `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
		Event     TransactionEvent   `json:"event"`
		Sender    string             `json:"sender"`
		Receiver  string             `json:"receiver"`
		Amount    float32            `json:"amount"`
		Currency  *Currency          `json:"currency"`
		Hash      string             `json:"hash"`
		Status    *TransactionStatus `json:"status"`
		Reason    string             `json:"reason"`
		CreatedAt time.Time          `json:"created_at"`
		UpdatedAt time.Time          `json:"updated_at"`
	}{
		Id:        transaction.Id,
		Status:    transaction.Status,
		Event:     PENDING,
		Sender:    transaction.Sender,
		Receiver:  transaction.Receiver,
		Amount:    transaction.Amount,
		Currency:  transaction.Currency,
		Hash:      transaction.Hash,
		CreatedAt: transaction.CreatedAt,
		UpdatedAt: transaction.UpdatedAt,
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
