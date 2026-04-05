package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/OtavMacedo/url-shortener-golang/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UrlRepository struct {
	pool *pgxpool.Pool
}

var ErrUrlSlugConflict = errors.New("url slug already exists")

func NewUrlRepository(pool *pgxpool.Pool) *UrlRepository {
	return &UrlRepository{pool: pool}
}

func (ur *UrlRepository) Create(ctx context.Context, url model.UrlModel) error {
	query := `INSERT INTO urls (id, user_id, original_url, slug) VALUES ($1, $2, $3, $4)`
	_, err := ur.pool.Exec(ctx, query,
		url.ID,
		url.UserID,
		url.OriginalUrl,
		url.Slug,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrUrlSlugConflict
		}
		return fmt.Errorf("failed to create url: %w", err)
	}
	return nil
}

func (ur *UrlRepository) FindBySlug(ctx context.Context, slug string) (*model.UrlModel, error) {
	query := `SELECT id, user_id, original_url, slug FROM urls WHERE slug = $1`

	var url model.UrlModel
	err := ur.pool.QueryRow(ctx, query, slug).Scan(
		&url.ID,
		&url.UserID,
		&url.OriginalUrl,
		&url.Slug,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find url by slug: %w", err)
	}
	return &url, nil
}
