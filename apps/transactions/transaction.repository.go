package transaction

import (
	"fmt"
	"myapp/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func createTransaction(db *gorm.DB, transaction *models.Transaction) error {
	return db.Create(transaction).Error
}

func updateTransaction(db *gorm.DB, transactionId uuid.UUID, updates models.Transaction) (models.Transaction, error) {
	var transaction models.Transaction

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
