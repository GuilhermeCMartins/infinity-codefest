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
