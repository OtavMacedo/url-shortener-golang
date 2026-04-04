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

type UserRepository struct {
	pool *pgxpool.Pool
}

var ErrUserEmailConflict = errors.New("user email already exists")

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (ur *UserRepository) Create(ctx context.Context, user model.UserModel) error {
	query := `
		INSERT INTO users (id, email, password_hash)
		VALUES ($1, $2, $3)
	`
	_, err := ur.pool.Exec(ctx, query,
		user.ID,
		user.Email,
		user.PasswordHash,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrUserEmailConflict
		}
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (ur *UserRepository) FindByEmail(ctx context.Context, email string) (*model.UserModel, error) {
	query := `
		SELECT id, email, password_hash FROM users WHERE email = $1
	`

	var user model.UserModel
	err := ur.pool.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}
	return &user, nil
}
