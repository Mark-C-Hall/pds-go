package service

import (
	"context"
	"fmt"
	"time"

	"github.com/bluesky-social/indigo/atproto/syntax"
	"github.com/mark-c-hall/pds-go/internal/core/model"
)

type AccountService interface {
	CreateAccount(ctx context.Context, handle syntax.Handle, email, password string) (*model.Account, error)
}

type AccountServiceImpl struct{}

func NewAccountService() AccountService {
	return &AccountServiceImpl{}
}

func (s *AccountServiceImpl) CreateAccount(ctx context.Context, handle syntax.Handle, email, password string) (*model.Account, error) {
	if handle == "" {
		return nil, fmt.Errorf("handle cannot be empty")
	}
	if email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}
	if password == "" {
		return nil, fmt.Errorf("password cannot be empty")
	}

	DID := fmt.Sprintf("did:plc:%s", handle)

	return &model.Account{
		DID:       syntax.DID(DID),
		Handle:    syntax.Handle(handle),
		Email:     email,
		CreatedAt: time.Now().UTC(),
	}, nil
}
