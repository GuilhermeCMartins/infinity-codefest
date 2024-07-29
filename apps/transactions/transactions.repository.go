package transactions

import (
	"fmt"
	"myapp/db"
	"myapp/models"

	"github.com/google/uuid"
)

func createTransaction(transaction *models.Transaction) error {
	db := db.Init()

	return db.Create(transaction).Error
}

func FindTxById(id uuid.UUID) (models.Transaction, error) {
	var tx models.Transaction
	db := db.Init()

	result := db.First(&tx, "id = ?", id)

	if result.Error != nil {
		return models.Transaction{}, result.Error
	}

	return tx, nil
}

func updateTransaction(transactionId uuid.UUID, updates models.Transaction) (models.Transaction, error) {
	var transaction models.Transaction
	db := db.Init()

	result := db.First(&transaction, "id = ?", transactionId)
	if result.Error != nil {
		return models.Transaction{}, fmt.Errorf("Transaction not found: %v", result.Error)
	}

	if err := db.Model(&transaction).Updates(updates).Error; err != nil {
		return models.Transaction{}, fmt.Errorf("Failed to update transaction: %v", err)
	}

	if err := db.First(&transaction, "id = ?", transactionId).Error; err != nil {
		return models.Transaction{}, fmt.Errorf("Failed to fetch updated transaction: %v", err)
	}

	return transaction, nil
}

func FindAllTransactions() (transactions []models.Transaction, count int, err error) {
	db := db.Init()

	result := db.Model(&models.Transaction{}).Preload("Users").Find(&transactions)

	if result.Error != nil {
		return nil, 0, result.Error
	}
	count = len(transactions)
	return transactions, count, nil
}
