package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
	Verify(hashedPassword, plainPassword string) error
}

type BcryptPasswordHasher struct {
	cost int
}

func NewBcryptPasswordHasher() PasswordHasher {
	return &BcryptPasswordHasher{
		cost: bcrypt.DefaultCost,
	}
}

func (h *BcryptPasswordHasher) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hash), nil
}

func (h *BcryptPasswordHasher) Verify(hashedPassword, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}
