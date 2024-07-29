package transactions

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func createTransaction(db *gorm.DB, transaction *Transaction) error {
	return db.Create(transaction).Error
}

func updateTransaction(db *gorm.DB, transactionId uuid.UUID, updates Transaction) (Transaction, error) {
	var transaction Transaction

	result := db.First(&transaction, "id = ?", transactionId)
	if result.Error != nil {
		return Transaction{}, fmt.Errorf("Transaction not found: %v", result.Error)
	}

	if err := db.Model(&transaction).Updates(updates).Error; err != nil {
		return Transaction{}, fmt.Errorf("Failed to update transaction: %v", err)
	}

	if err := db.First(&transaction, "id = ?", transactionId).Error; err != nil {
		return Transaction{}, fmt.Errorf("Failed to fetch updated transaction: %v", err)
	}

	return transaction, nil
}
