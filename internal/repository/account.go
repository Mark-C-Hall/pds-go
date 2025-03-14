package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mark-c-hall/pds-go/internal/model"
)

type AccountRepository interface {
	CreateAccount(ctx context.Context, account *model.Account, hashedPassword string) error
}

type SQLAccountRepository struct {
	db *sql.DB
}

func NewSQLAccountRepository(db *sql.DB) AccountRepository {
	return &SQLAccountRepository{db: db}
}

func (r *SQLAccountRepository) CreateAccount(ctx context.Context, account *model.Account, hashedPassword string) error {
	query := `
		INSERT INTO accounts (did, handle, email, password_hash, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.ExecContext(ctx, query,
		account.DID.String(),
		account.Handle.String(),
		account.Email,
		hashedPassword,
		account.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create account: %w", err)
	}

	return nil
}
