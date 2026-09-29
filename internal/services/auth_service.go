package services

import (
	"context"
	"drivex/internal/models"
	"drivex/internal/repository"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(
	userRepo *repository.UserRepository,
) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}


func (s *AuthService) Register(
	ctx context.Context,
	name string,
	email string,
	password string,
) (*models.User, error) {

	if name == "" {
		return nil, errors.New("name is required")
	}

	if email == "" {
		return nil, errors.New("email is required")
	}

	if password == "" {
		return nil, errors.New("password is required")
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.Create(
		ctx,
		name,
		email,
		string(passwordHash),
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}