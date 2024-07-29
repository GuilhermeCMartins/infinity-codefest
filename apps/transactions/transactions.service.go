package transactions

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"myapp/apps/user"
	"myapp/models"
	"time"

	"myapp/utils"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func verifyIfCreationIsValid(payload models.TransactionPayload) error {
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

func verifySignature(publicKey, signature string, sender, receiver uuid.UUID, amount float64, createdAt time.Time, currency models.Currency) (bool, error) {
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

func handleRequestTransaction(db *gorm.DB, payload models.TransactionPayload) string {
	error := verifyIfCreationIsValid(payload)
	if error != nil {
		fmt.Errorf("[TRANSACTION]: Transaction request: %v", error)
		return ""
	}

	var sender models.User
	if err := db.First(&sender, "id = ?", payload.Sender).Error; err != nil {
		fmt.Errorf("[TRANSACTION]: User not found: %v", err)
		return ""
	}

	var receiver models.User
	if err := db.First(&receiver, "id = ?", payload.Receiver).Error; err != nil {
		fmt.Errorf("[TRANSACTION]: User not found: %v", err)
		return ""
	}

	if sender.Id == receiver.Id {
		fmt.Errorf("[TRANSACTION]: Cannot perform transaction between the same person")
		return ""
	}

	valid, err := verifySignature(sender.PublicKey, payload.Hash, sender.Id, receiver.Id, float64(payload.Amount), payload.CreatedAt, *payload.Currency)
	if err != nil {
		fmt.Errorf("[TRANSACTION]: Erro ao decodificar assinatura: %v", err)
		return ""
	}

	if !valid {
		fmt.Errorf("[TRANSACTION]: Invalid hash")
		return ""
	}

	_, amountInSenderCurrency, errConvert := utils.ConvertCurrency(float64(payload.Amount), utils.Currency(*sender.Currency), utils.Currency(*receiver.Currency), utils.Currency(*payload.Currency))
	if errConvert != nil {
		fmt.Printf("[TRANSACTION]: Conversion failed: %v", errConvert)
		return ""
	}

	hasValue := float64(sender.Balance) - amountInSenderCurrency

	if hasValue < 0 {
		fmt.Errorf("[TRANSACTION]: Insuficient balance")
		return ""
	}

	status := models.TX_REVIEW

	transaction := models.Transaction{
		Id:        uuid.New(),
		Status:    &status,
		Sender:    payload.Sender,
		Receiver:  payload.Receiver,
		Amount:    payload.Amount,
		Currency:  payload.Currency,
		Hash:      payload.Hash,
		Users:     []*models.User{&sender, &receiver},
		CreatedAt: payload.CreatedAt,
		UpdatedAt: payload.UpdatedAt,
	}

	if err := createTransaction(&transaction); err != nil {
		fmt.Printf("[TRANSACTION] Failed to create transaction: %v", err)
		return ""
	}

	fmt.Println("[CREATED TRANSACTION]")

	message := struct {
		Id        uuid.UUID                 `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
		Event     models.TransactionEvent   `json:"event"`
		Sender    uuid.UUID                 `json:"sender" gorm:"type:uuid"`
		Receiver  uuid.UUID                 `json:"receiver" gorm:"type:uuid"`
		Amount    float32                   `json:"amount"`
		Currency  *models.Currency          `json:"currency"`
		Hash      string                    `json:"hash"`
		Status    *models.TransactionStatus `json:"status"`
		Reason    string                    `json:"reason"`
		CreatedAt time.Time                 `json:"created_at"`
		UpdatedAt time.Time                 `json:"updated_at"`
	}{
		Id:        transaction.Id,
		Status:    transaction.Status,
		Event:     models.TX_PENDING,
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

func handlePendingTransaction(db *gorm.DB, payload models.TransactionPayload) string {
	error := verifyIfCreationIsValid(payload)
	if error != nil {
		fmt.Errorf("[TRANSACTION]: Transaction request: %v", error)
		return ""
	}

	var sender models.User
	if err := db.First(&sender, "id = ?", payload.Sender).Error; err != nil {
		fmt.Errorf("[TRANSACTION]: User not found: %v", err)
		return ""
	}

	var receiver models.User
	if err := db.First(&receiver, "id = ?", payload.Receiver).Error; err != nil {
		fmt.Errorf("[TRANSACTION]: User not found: %v", err)
		return ""
	}

	var transaction models.Transaction
	result := db.First(&transaction, "id = ?", payload.Id)
	if result.Error != nil {
		fmt.Errorf("[TRANSACTION]: Transaction not found: %v", result.Error)
		return ""
	}

	status := models.TX_APPROVED

	if *payload.Status == models.TX_FAILED {
		status := models.TX_FAILED

		updates := models.Transaction{
			Status:    &status,
			Reason:    "Reprovado pelo KYC/FRAUD",
			UpdatedAt: time.Now(),
		}

		transactionUpdated, err := updateTransaction(transaction.Id, updates)
		if err != nil {
			fmt.Errorf("[TRANSACTION]: Failed to update transaction: %v", err)
			return ""
		}

		message := utils.CreateMessage(transactionUpdated, models.TX_PENDING)
		return message
	}

	updates := models.Transaction{
		Status:    &status,
		Reason:    payload.Reason,
		UpdatedAt: time.Now(),
	}

	amountInReiceverCurrency, amountInSenderCurrency, errConvert := utils.ConvertCurrency(float64(payload.Amount), utils.Currency(*sender.Currency), utils.Currency(*receiver.Currency), utils.Currency(*payload.Currency))
	if errConvert != nil {
		fmt.Errorf("[TRANSACTION]: Conversion failed: %v", errConvert)
		return ""
	}

	newValueReceiver := amountInReiceverCurrency + receiver.Balance
	newValueSender := amountInSenderCurrency - sender.Balance

	updateSender := models.User{
		Balance: newValueSender,
	}
	updateReceiver := models.User{
		Balance: newValueReceiver,
	}

	transactionUpdated, err := updateTransaction(transaction.Id, updates)
	if err != nil {
		fmt.Errorf("[TRANSACTION]: Failed to update transaction: %v", err)
		return ""
	}

	_, err = user.UpdateUser(receiver.Id, updateReceiver)
	if err != nil {
		fmt.Errorf("[TRANSACTION]: Failed to update receiver user: %v", err)
		return ""
	}

	_, err = user.UpdateUser(sender.Id, updateSender)
	if err != nil {
		fmt.Errorf("[TRANSACTION]: Failed to update sender user: %v", err)
		return ""
	}

	message := utils.CreateMessage(transactionUpdated, models.TX_CREATED)
	fmt.Println("[TRANSACTION PENDING]:", message)
	return message
}

func HandleMessageTransaction(db *gorm.DB, payload models.TransactionPayload) string {

	var result string

	switch payload.Event {
	case models.TX_REQUEST:
		result = handleRequestTransaction(db, payload)
	case models.TX_PENDING:
		result = handlePendingTransaction(db, payload)
	default:
		fmt.Printf("Unknown event type: %s", payload.Event)
	}

	return result
}
