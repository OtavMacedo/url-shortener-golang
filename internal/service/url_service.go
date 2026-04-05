package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"

	"github.com/OtavMacedo/url-shortener-golang/internal/model"
	"github.com/OtavMacedo/url-shortener-golang/internal/repository"
	"github.com/OtavMacedo/url-shortener-golang/utils"
)

type UrlService struct {
	repository *repository.UrlRepository
}

func NewUrlService(repository *repository.UrlRepository) *UrlService {
	return &UrlService{repository: repository}
}

func (us *UrlService) Create(ctx context.Context, url model.UrlModel) error {
	if url.Slug == "" {
		generated, err := generateSlug(6)
		if err != nil {
			return fmt.Errorf("failed to generate slug: %w", err)
		}
		url.Slug = generated
	}
	existingSlug, err := us.repository.FindBySlug(ctx, url.Slug)
	if err != nil {
		return fmt.Errorf("failed to verify existing url: %w", err)
	}
	if existingSlug != nil {
		return repository.ErrUrlSlugConflict
	}
	url.ID, err = utils.NewUUIDv7()
	if err != nil {
		return fmt.Errorf("failed to generate url id: %w", err)
	}
	if err := us.repository.Create(ctx, url); err != nil {
		if errors.Is(err, repository.ErrUrlSlugConflict) {
			return repository.ErrUrlSlugConflict
		}
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func generateSlug(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		b[i] = charset[n.Int64()]
	}
	return string(b), nil
}
