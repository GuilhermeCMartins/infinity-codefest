package user

import (
	"fmt"
	"myapp/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateUser(db *gorm.DB, user *models.User) error {
	return db.Create(user).Error
}

func UpdateUser(db *gorm.DB, userId uuid.UUID, updates models.User) (models.User, error) {
	var user models.User

	result := db.First(&user, "id = ?", userId)
	if result.Error != nil {
		return models.User{}, fmt.Errorf("User not found: %v", result.Error)
	}

	if err := db.Model(&user).Updates(updates).Error; err != nil {
		return models.User{}, fmt.Errorf("Failed to update user: %v", err)
	}

	if err := db.First(&user, "id = ?", userId).Error; err != nil {
		return models.User{}, fmt.Errorf("Failed to fetch updated user: %v", err)
	}

	return user, nil
}

func FindAllUsers(db *gorm.DB) (users []models.User, count int, err error) {
	result := db.Find(&users)
	
	if result.Error != nil {
			return nil, 0, result.Error
	}
	
	count = len(users)
	
	return users, count, nil
}

func FindUserById(db *gorm.DB, id uuid.UUID) (models.User, error) {
	var user models.User
	result := db.First(&user, "id = ?", id)
	
	if result.Error != nil {
		return models.User{}, result.Error
	}
	
	return user, nil
}

func FindUserTransactions(db *gorm.DB, userID uuid.UUID) ([]models.Transaction, int, error) {
	var transactions []models.Transaction

	result := db.Where("sender = ? OR receiver = ?", userID, userID).Find(&transactions)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	count := len(transactions)

	return transactions, count, nil
}

func FindUserTransactionByTransactionId(db *gorm.DB, userID uuid.UUID, txID uuid.UUID) (models.Transaction, string, error) {
	var transaction models.Transaction

	result := db.Where("sender = ? OR receiver = ?", userID, userID).First(&transaction, "id = ?", txID)

	sender := transaction.Sender.String()

	if result.Error != nil {
		return models.Transaction{}, "", result.Error
	}

	return transaction, sender, nil
}

func FindUserTransactionsByStatus(db *gorm.DB, userID uuid.UUID, status models.TransactionStatus) ([]models.Transaction, int, error) {
	var transactions []models.Transaction

	result := db.Where("sender = ? OR receiver = ?", userID, userID).Where("status LIKE ?", status).Find(&transactions)

	if result.Error != nil {
		return nil, 0, result.Error
	}
	count := len(transactions)

	return transactions, count, nil
}
