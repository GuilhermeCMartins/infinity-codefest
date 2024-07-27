package user

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateUser(db *gorm.DB, user *User) error {
	return db.Create(user).Error
}

func updateUser(db *gorm.DB, userId uuid.UUID, updates User) (User, error) {
	var user User

	result := db.First(&user, "id = ?", userId)
	if result.Error != nil {
		return User{}, fmt.Errorf("User not found: %v", result.Error)
	}

	if err := db.Model(&user).Updates(updates).Error; err != nil {
		return User{}, fmt.Errorf("Failed to update user: %v", err)
	}

	if err := db.First(&user, "id = ?", userId).Error; err != nil {
		return User{}, fmt.Errorf("Failed to fetch updated user: %v", err)
	}

	return user, nil
}
