package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/OtavMacedo/url-shortener-golang/internal/model"
	"github.com/OtavMacedo/url-shortener-golang/internal/repository"
	"github.com/OtavMacedo/url-shortener-golang/utils"
)

var ErrUserAlreadyExists = errors.New("user already exists")

type UserService struct {
	repository *repository.UserRepository
}

func NewUserService(repository *repository.UserRepository) *UserService {
	return &UserService{repository: repository}
}

func (us *UserService) Create(ctx context.Context, user model.UserModel) error {
	existingUser, err := us.repository.FindByEmail(ctx, user.Email)
	if err != nil {
		return fmt.Errorf("failed to verify existing user: %w", err)
	}

	if existingUser != nil {
		return ErrUserAlreadyExists
	}

	user.ID, err = utils.NewUUIDv7()
	if err != nil {
		return fmt.Errorf("failed to generate user id: %w", err)
	}

	if err := us.repository.Create(ctx, user); err != nil {
		if errors.Is(err, repository.ErrUserEmailConflict) {
			return ErrUserAlreadyExists
		}

		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}
