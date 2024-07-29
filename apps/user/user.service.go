package user

import (
	"encoding/json"
	"fmt"
	"log"
	"myapp/models"
	"myapp/utils"
	"time"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

func verifyIfCreationIsValid(payload models.UserPayload) error {
	validate := validator.New()
	err := validate.Struct(payload)
	if err != nil {
		return err
	}
	return nil
}

func handleRequestUser(payload models.UserPayload) string {
	err := verifyIfCreationIsValid(payload)

	status := models.USER_REVIEW

	user := models.User{
		Id:        uuid.New(),
		Name:      payload.Name,
		Email:     payload.Email,
		PublicKey: payload.PublicKey,
		Status:    &status,
		Balance:   payload.Balance,
		Currency:  &payload.Currency,
		CreatedAt: payload.CreatedAt,
		UpdatedAt: payload.UpdatedAt,
	}

	if err != nil {
		failedStatus := models.USER_FAILED

		updates := models.User{
			Status:    &failedStatus,
			Reason:    "Falta de campos para criação de usuário",
			UpdatedAt: time.Now(),
		}

		userUpdated, _ := UpdateUser(user.Id, updates)

		message := utils.CreateMessage(userUpdated, models.USER_PENDING)

		return message
	}

	err = CreateUser(&user)
	if err != nil {
		fmt.Printf("Failed to create user: %v", err)
		return ""
	}

	message := struct {
		Id        uuid.UUID          `json:"id"`
		Status    *models.UserStatus `json:"status"`
		Event     models.UserEvents  `json:"event" validate:"required"`
		Name      string             `json:"name" validate:"required"`
		Email     string             `json:"email" validate:"required,email"`
		PublicKey string             `json:"public_key" validate:"required"`
		Balance   float64            `json:"balance" validate:"required"`
		Currency  models.Currency    `json:"currency" validate:"required"`
		CreatedAt time.Time          `json:"created_at" validate:"required"`
		UpdatedAt time.Time          `json:"updated_at"`
	}{
		Id:        user.Id,
		Status:    user.Status,
		Event:     models.USER_PENDING,
		Name:      payload.Name,
		Email:     payload.Email,
		PublicKey: payload.PublicKey,
		Balance:   payload.Balance,
		Currency:  payload.Currency,
		CreatedAt: payload.CreatedAt,
		UpdatedAt: payload.UpdatedAt,
	}

	messageJSON, err := json.Marshal(message)

	if err != nil {
		panic(err)
	}

	var stringMessage = string(messageJSON)

	return stringMessage
}

func handlePendingUser(payload models.UserPayload) string {
	err := verifyIfCreationIsValid(payload)
	if err != nil {
		log.Printf("Validation failed: %v", err)
		return ""
	}

	result, err := FindUserById(payload.Id)
	if err != nil {
		log.Printf("User not found: %v", err)
		return ""
	}

	status := models.USER_APPROVED

	if *payload.Status == models.USER_FAILED {
		status := models.USER_FAILED

		updates := models.User{
			Status:    &status,
			Reason:    "Reprovado pelo KYC/FRAUD",
			UpdatedAt: time.Now(),
		}

		userUpdated, err := UpdateUser(result.Id, updates)
		if err != nil {
			fmt.Printf("[USER]: Error on update user: %v", err)
			return ""
		}

		message := utils.CreateMessage(userUpdated, models.USER_PENDING)
		return message
	}

	updates := models.User{
		Status:    &status,
		Reason:    payload.Reason,
		UpdatedAt: time.Now(),
	}

	userUpdated, err := UpdateUser(result.Id, updates)
	if err != nil {
		fmt.Printf("[USER]: Error on update user: %v", err)
		return ""
	}

	message := utils.CreateMessage(userUpdated, models.USER_CREATED)
	println("[USER PENDING]", message)
	return message
}

func HandleMessageUser(payload models.UserPayload) string {
	var result string

	switch payload.Event {
	case models.USER_REQUEST:
		result = handleRequestUser(payload)
	case models.USER_PENDING:
		result = handlePendingUser(payload)
	default:
		log.Printf("Unknown event type: %s", payload.Event)
	}

	return result
}

type UsersResponse struct {
	Users []models.User `json:"users"`
	Count int           `json:"count"`
}
