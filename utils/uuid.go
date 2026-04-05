package utils

import (
	"fmt"

	"github.com/google/uuid"
)

func NewUUIDv7() (string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return "", fmt.Errorf("failed to generate id: %w", err)
	}
	return id.String(), nil
}
