package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/OtavMacedo/url-shortener-golang/internal/apperr"
	"github.com/OtavMacedo/url-shortener-golang/internal/cache"
	"github.com/OtavMacedo/url-shortener-golang/internal/model"
	"github.com/OtavMacedo/url-shortener-golang/internal/repository"
	"github.com/OtavMacedo/url-shortener-golang/utils"
)

type UrlService struct {
	repository *repository.UrlRepository
	cache      *cache.RedisCache
}

func NewUrlService(repository *repository.UrlRepository, cache *cache.RedisCache) *UrlService {
	return &UrlService{repository: repository, cache: cache}
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

func (us *UrlService) FindBySlug(ctx context.Context, slug string) (string, error) {
	originalUrl, err := us.cache.Get(ctx, "url:"+slug)
	if err == nil && originalUrl != "" {
		fmt.Println("Retriving from cache")
		return originalUrl, nil
	}

	existingSlug, err := us.repository.FindBySlug(ctx, slug)
	if err != nil {
		return "", err
	}
	if existingSlug == nil {
		return "", apperr.ErrSlugNotFound
	}
	us.cache.Set(ctx, "url:"+slug, existingSlug.OriginalUrl, 24*time.Hour)
	return existingSlug.OriginalUrl, nil
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
