package transactions

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"myapp/apps/user"
	"myapp/models"
	"time"

	"myapp/utils"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
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

func handleRequestTransaction(payload models.TransactionPayload) string {
	error := verifyIfCreationIsValid(payload)
	if error != nil {
		fmt.Printf("[TRANSACTION]: Transaction request: %v", error)
		return ""
	}

	if payload.Receiver == payload.Sender {
		fmt.Printf("[TRANSACTION]: Transaction for the same person")
		return ""
	}

	sender, err := user.FindUserById(payload.Sender)
	if err != nil {
		fmt.Printf("[TRANSACTION]: Sender not found: %v", err)
		return ""
	}

	receiver, err := user.FindUserById(payload.Receiver)
	if err != nil {
		fmt.Printf("[TRANSACTION]: Receiver not found: %v", err)
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

	hasValue := sender.Balance - amountInSenderCurrency

	if hasValue < 0 {
		fmt.Println("[TRANSACTION]: Insuficient balance")
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

func handlePendingTransaction(payload models.TransactionPayload) string {
	if payload.Event == models.TX_CREATED {
		fmt.Printf("[TRANSACTION]: Already created")
		return ""
	}

	error := verifyIfCreationIsValid(payload)
	if error != nil {
		fmt.Printf("[TRANSACTION]: Transaction request: %v", error)
		return ""
	}

	if payload.Receiver == payload.Sender {
		fmt.Printf("[TRANSACTION]: Transaction for the same person")
		return ""
	}

	sender, err := user.FindUserById(payload.Sender)
	if err != nil {
		fmt.Printf("[TRANSACTION]: Sender not found: %v", err)
		return ""
	}

	userFailed := models.USER_FAILED

	if sender.Status == &userFailed {
		fmt.Printf("[TRANSACTION]: Sender status failed")
		return ""
	}

	receiver, err := user.FindUserById(payload.Receiver)
	if err != nil {
		fmt.Printf("[TRANSACTION]: Sender not found: %v", err)
		return ""
	}

	if receiver.Status == &userFailed {
		fmt.Printf("[TRANSACTION]: Receiver status failed")
		return ""
	}

	transaction, err := FindTxById(payload.Id)
	if err != nil {
		fmt.Printf("[TRANSACTION]: Transaction not found: %v", err)
		return ""
	}

	statusCompleted := models.TX_APPROVED
	statusFailed := models.TX_FAILED
	if transaction.Status == &statusCompleted || transaction.Status == &statusFailed {
		fmt.Printf("[TRANSACTION]: Already processed")
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

	amountInReiceverCurrency, amountInSenderCurrency, errConvert := utils.ConvertCurrency(float64(payload.Amount), utils.Currency(*sender.Currency), utils.Currency(*receiver.Currency), utils.Currency(*payload.Currency))
	if errConvert != nil {
		fmt.Printf("[TRANSACTION]: Conversion failed: %v", errConvert)
		return ""
	}

	balance := utils.Float64ToBigFloat(sender.Balance)
	amount := utils.Float64ToBigFloat(amountInSenderCurrency)

	result := new(big.Float).Sub(balance, amount)

	if result.Cmp(big.NewFloat(0)) < 0 {
		fmt.Println("Insufficient balance.")
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
			fmt.Printf("[TRANSACTION]: Error on update transaction: %v", err)
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

	balanceReceiver := utils.Float64ToBigFloat(receiver.Balance)
	amountReceiver := utils.Float64ToBigFloat(amountInReiceverCurrency)
	balanceSender := utils.Float64ToBigFloat(sender.Balance)
	amountSender := utils.Float64ToBigFloat(amountInSenderCurrency)

	newValueReceiver := new(big.Float).Add(balanceReceiver, amountReceiver)
	newValueSender := new(big.Float).Sub(balanceSender, amountSender)

	newReceiverBalance, _ := utils.BigFloatToFloat64(newValueReceiver)

	newSenderBalance, _ := utils.BigFloatToFloat64(newValueSender)

	updateSender := models.User{
		Balance: newSenderBalance,
	}
	updateReceiver := models.User{
		Balance: newReceiverBalance,
	}

	transactionUpdated, _ := updateTransaction(transaction.Id, updates)
	user.UpdateUser(receiver.Id, updateReceiver)
	user.UpdateUser(sender.Id, updateSender)

	message := utils.CreateMessage(transactionUpdated, models.TX_CREATED)
	fmt.Println("[TRANSACTION COMPLETED]:", message)
	return message
}

func HandleMessageTransaction(payload models.TransactionPayload) string {

	var result string

	switch payload.Event {
	case models.TX_REQUEST:
		result = handleRequestTransaction(payload)
	case models.TX_PENDING:
		result = handlePendingTransaction(payload)
	default:
		fmt.Printf("Unknown event type: %s", payload.Event)
	}

	return result
}
