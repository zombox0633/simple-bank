package user

import (
	"context"
	"fmt"

	db "simplebank/db/sqlc"
	"simplebank/internal/common"
)

// Service owns user business rules independently from Gin and HTTP binding.
type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func (service *Service) CreateUser(
	ctx context.Context,
	username string,
	password string,
	fullName string,
	email string,
) (db.User, error) {
	hashedPassword, err := common.HashPassword(password)
	if err != nil {
		return db.User{}, fmt.Errorf("hash password: %w", err)
	}

	createdUser, err := service.store.CreateUser(ctx, db.CreateUserParams{
		Username: username,
		Password: hashedPassword,
		FullName: fullName,
		Email:    email,
	})
	if err != nil {
		if db.SQLState(err) == db.SQLStateUniqueViolation {
			return db.User{}, common.ErrConflict
		}
		return db.User{}, fmt.Errorf("create user: %w", err)
	}

	return createdUser, nil
}
