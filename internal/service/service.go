package service

import (
	"context"
	"errors"
	"github.com/castiglionimax/PocEventSourcingAccounting/internal/domain"
)

type (
	repository interface {
		CreateAccount(ctx context.Context, account domain.Account) (domain.AccountID, error)
		GetBalance(ctx context.Context, accountID domain.AccountID) (float32, error)
		SaveTransaction(ctx context.Context, transaction domain.Transaction) error
	}

	Service struct {
		repository repository
	}
)

func NewService(repository repository) (*Service, error) {
	if repository == nil {
		return nil, errors.New("repository should not be nil")
	}
	return &Service{repository: repository}, nil
}

func (s Service) CreateAccount(ctx context.Context, account domain.Account) (domain.AccountID, error) {
	return s.repository.CreateAccount(ctx, account)
}

func (s Service) Transaction(ctx context.Context, transaction domain.Transaction) error {
	return s.repository.SaveTransaction(ctx, transaction)
}

func (s Service) GetBalance(ctx context.Context, accountID domain.AccountID) (float32, error) {
	return s.repository.GetBalance(ctx, accountID)
}
