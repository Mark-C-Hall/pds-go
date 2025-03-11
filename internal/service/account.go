package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bluesky-social/indigo/atproto/syntax"
	"github.com/mark-c-hall/pds-go/internal/model"
	"github.com/mark-c-hall/pds-go/internal/repository"
	"github.com/mark-c-hall/pds-go/internal/util"
)

type AccountService interface {
	CreateAccount(ctx context.Context, handle syntax.Handle, email, password string) (*model.Account, error)
}

type AccountServiceImpl struct {
	repo     repository.AccountRepository
	pwHasher util.PasswordHasher
	logger   *log.Logger
}

func NewAccountService(repo repository.AccountRepository, pwHasher util.PasswordHasher, logger *log.Logger) AccountService {
	return &AccountServiceImpl{
		repo:     repo,
		pwHasher: pwHasher,
		logger:   logger,
	}
}

func (s *AccountServiceImpl) CreateAccount(ctx context.Context, handle syntax.Handle, email, password string) (*model.Account, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	if handle == "" {
		return nil, fmt.Errorf("handle cannot be empty")
	}
	if email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}
	if password == "" {
		return nil, fmt.Errorf("password cannot be empty")
	}

	// TODO: check for unique handle

	DID := fmt.Sprintf("did:plc:%s", handle)

	hashedPassword, err := s.pwHasher.Hash(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	account := &model.Account{
		DID:       syntax.DID(DID),
		Handle:    syntax.Handle(handle),
		Email:     email,
		CreatedAt: time.Now().UTC(),
	}

	if err := s.repo.CreateAccount(ctx, account, hashedPassword); err != nil {
		return nil, fmt.Errorf("failed to store account: %w", err)
	}

	s.logger.Printf("Created account for %s with DID %s", account.Handle, account.DID)

	return account, nil
}
