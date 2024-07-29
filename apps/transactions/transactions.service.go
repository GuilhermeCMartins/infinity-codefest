package transactions

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"myapp/apps/user"
	"time"

	"myapp/utils"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func createMessage(transaction Transaction, event TransactionEvent) string {
	message := struct {
		Id        uuid.UUID          `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
		Sender    string             `json:"sender"`
		Receiver  string             `json:"receiver"`
		Amount    float32            `json:"amount"`
		Currency  Currency           `json:"currency"`
		Hash      string             `json:"hash"`
		Status    *TransactionStatus `json:"status"`
		Reason    string             `json:"reason"`
		CreatedAt time.Time          `json:"created_at"`
		UpdatedAt time.Time          `json:"updated_at"`
		Event     TransactionEvent   `json:"event" validate:"required"`
	}{
		Id:        transaction.Id,
		Status:    transaction.Status,
		Sender:    transaction.Sender,
		Receiver:  transaction.Receiver,
		Hash:      transaction.Hash,
		Reason:    transaction.Reason,
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

func decodePublicKey(pubKeyHex string) (*ecdsa.PublicKey, error) {
	pubKeyBytes, err := hex.DecodeString(pubKeyHex)
	if err != nil {
		return nil, err
	}

	if len(pubKeyBytes) != 65 {
		return nil, fmt.Errorf("invalid public key length: %d", len(pubKeyBytes))
	}

	pubKey, err := crypto.UnmarshalPubkey(pubKeyBytes)
	if err != nil {
		return nil, err
	}

	return pubKey, nil
}

func floatToBytes(amount float64) []byte {
	return []byte(fmt.Sprintf("%.2f", amount))
}

func verifySignature(publicKey, signature string, sender, receiver uuid.UUID, amount float64, createdAt time.Time, currency Currency) (bool, error) {
	senderBytes := sender[:]
	receiverBytes := receiver[:]
	amountBytes := floatToBytes(amount)
	createdAtBytes := []byte(createdAt.Format(time.RFC3339))
	currencyBytes := []byte(currency)

	data := append(senderBytes, receiverBytes...)
	data = append(data, amountBytes...)
	data = append(data, createdAtBytes...)
	data = append(data, currencyBytes...)

	hash := crypto.Keccak256Hash(data)

	pubKey, _ := decodePublicKey(publicKey)

	publicKeyBytes := crypto.FromECDSAPub(pubKey)

	signatureBytes, _ := hex.DecodeString(signature)

	sigPublicKey, _ := crypto.Ecrecover(hash.Bytes(), signatureBytes)

	matches := bytes.Equal(sigPublicKey, publicKeyBytes)

	return matches, nil
}

// TO-DO: Não permitir que faca transacao entre a mesma pessoa
func handleRequestTransaction(db *gorm.DB, payload TransactionPayload) string {
	error := verifyIfCreationIsValid(payload)
	if error != nil {
		fmt.Printf("[TRANSACTION]: Transaction request: %v", error)
		return ""
	}

	var sender user.User
	if err := db.First(&sender, "id = ?", payload.Sender).Error; err != nil {
		fmt.Printf("[TRANSACTION]: User not found: %v", err)
		return ""
	}

	var receiver user.User
	if err := db.First(&receiver, "id = ?", payload.Receiver).Error; err != nil {
		fmt.Printf("[TRANSACTION]: User not found: %v", err)
		return ""
	}

	valid, err := verifySignature(sender.PublicKey, payload.Hash, sender.Id, receiver.Id, float64(payload.Amount), payload.CreatedAt, *payload.Currency)
	if err != nil {
		fmt.Errorf("[TRANSACTION]: Erro ao decodificar assinatura: %v", err)
		return ""
	}

	if !valid {
		fmt.Println("[TRANSACTION]: Invalid hash")
		return ""
	}

	_, amountInSenderCurrency, errConvert := utils.ConvertCurrency(float64(payload.Amount), utils.Currency(*sender.Currency), utils.Currency(*receiver.Currency), utils.Currency(*payload.Currency))
	if errConvert != nil {
		fmt.Printf("[TRANSACTION]: Conversion failed: %v", errConvert)
		return ""
	}

	hasValue := float64(sender.Balance) - amountInSenderCurrency

	if hasValue < 0 {
		fmt.Println("[TRANSACTION]: Insuficient balance")
		return ""
	}

	status := REVIEW

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

	if err := createTransaction(db, &transaction); err != nil {
		fmt.Printf("[TRANSACTION] Failed to create transaction: %v", err)
		return ""
	}

	fmt.Println("[CREATED TRANSACTION]")

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

	fmt.Println("[TRANSACTION SENDING PEDING]")

	return stringMessage
}

// TO-DO: Arrumar os updates e os retornos de mensagens
func handlePendingTransaction(db *gorm.DB, payload TransactionPayload) string {
	error := verifyIfCreationIsValid(payload)
	if error != nil {
		fmt.Printf("[TRANSACTION]: Transaction request: %v", error)
		return ""
	}

	var sender user.User
	if err := db.First(&sender, "id = ?", payload.Sender).Error; err != nil {
		fmt.Printf("[TRANSACTION]: User not found: %v", err)
		return ""
	}

	var receiver user.User
	if err := db.First(&receiver, "id = ?", payload.Receiver).Error; err != nil {
		fmt.Printf("[TRANSACTION]: User not found: %v", err)
		return ""
	}

	var transaction Transaction
	result := db.First(&transaction, "id = ?", payload.Id)
	if result.Error != nil {
		fmt.Printf("[TRANSACTION]: Transaction not found: %v", result.Error)
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
		//TO-DO: tratar erro de banco

		message := createMessage(transactionUpdated, PENDING)
		return message
	}

	updates := Transaction{
		Status:    &status,
		Reason:    payload.Reason,
		UpdatedAt: time.Now(),
	}

	amountInReiceverCurrency, amountInSenderCurrency, errConvert := utils.ConvertCurrency(float64(payload.Amount), utils.Currency(*sender.Currency), utils.Currency(*receiver.Currency), utils.Currency(*payload.Currency))
	if errConvert != nil {
		fmt.Printf("[TRANSACTION]: Conversion failed: %v", errConvert)
		return ""
	}

	newValueReceiver := amountInReiceverCurrency + receiver.Balance
	newValueSender := amountInSenderCurrency - sender.Balance

	updateSender := user.User{
		Balance: newValueSender,
	}
	updateReceiver := user.User{
		Balance: newValueReceiver,
	}

	transactionUpdated, _ := updateTransaction(db, transaction.Id, updates)
	user.UpdateUser(db, receiver.Id, updateReceiver)
	user.UpdateUser(db, sender.Id, updateSender)

	message := createMessage(transactionUpdated, CREATED)
	fmt.Println("[TRANSACTION PENDING]:", message)
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
		fmt.Printf("Unknown event type: %s", payload.Event)
	}

	return result
}

func FindAllTransactions(db *gorm.DB) (transactions []Transaction, count int, err error) {
	result := db.Model(&Transaction{}).Preload("User").Find(&transactions)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	count = len(transactions)
	return transactions, count, nil
}