package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/OtavMacedo/url-shortener-golang/internal/apperr"
	"github.com/OtavMacedo/url-shortener-golang/internal/dto"
	"github.com/OtavMacedo/url-shortener-golang/internal/model"
	"github.com/OtavMacedo/url-shortener-golang/internal/repository"
	"github.com/OtavMacedo/url-shortener-golang/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repository  *repository.UserRepository
	authService *AuthService
}

func NewUserService(repository *repository.UserRepository, authService *AuthService) *UserService {
	return &UserService{
		repository:  repository,
		authService: authService,
	}
}

func (us *UserService) Create(ctx context.Context, user dto.CreateUserInput) error {
	existingUser, err := us.repository.FindByEmail(ctx, user.Email)
	if err != nil {
		return fmt.Errorf("failed to verify existing user: %w", err)
	}

	if existingUser != nil {
		return apperr.ErrUserAlreadyExists
	}

	userID, err := utils.NewUUIDv7()
	if err != nil {
		return fmt.Errorf("failed to generate user id: %w", err)
	}
	hashedPassword, err := us.generatePasswordHash(user.Password)
	if err != nil {
		return err
	}

	userDb := model.UserModel{
		ID:           userID,
		Email:        user.Email,
		PasswordHash: string(hashedPassword),
	}

	if err := us.repository.Create(ctx, userDb); err != nil {
		if errors.Is(err, repository.ErrUserEmailConflict) {
			return apperr.ErrUserAlreadyExists
		}

		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (us *UserService) Login(ctx context.Context, login dto.LoginInput) (string, error) {
	existingUser, err := us.repository.FindByEmail(ctx, login.Email)
	if err != nil {
		return "", fmt.Errorf("failed to verify existing user: %w", err)
	}
	if existingUser == nil {
		return "", apperr.ErrUserNotFound
	}
	isCorrectPassword := us.comparePasswordHash(existingUser.PasswordHash, login.Password)
	if !isCorrectPassword {
		return "", apperr.ErrInvalidPassword
	}
	token, err := us.authService.GenerateToken(existingUser.ID)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (us *UserService) generatePasswordHash(password string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	return hash, nil
}

func (us *UserService) comparePasswordHash(hashedPassword, cleanPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(cleanPassword))
	if err != nil {
		return false
	}
	return true
}
