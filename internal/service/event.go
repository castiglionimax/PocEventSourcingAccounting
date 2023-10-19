package service

import (
	"context"

	"github.com/castiglionimax/PocEventSourcingAccounting/internal/domain"
)

type (
	projection interface {
		CreateAccount(ctx context.Context, account domain.Account) error
		RegisterTransaction(ctx context.Context, tx domain.Transaction) error
	}

	EventService struct {
		repository projection
	}
)

const deposit = "deposit"

func NewEventService(projection projection) *EventService {
	return &EventService{repository: projection}
}

func (e EventService) CreateAccount(ctx context.Context, account domain.Account) error {
	return e.repository.CreateAccount(ctx, account)
}

func (e EventService) RegisterTransaction(ctx context.Context, tx domain.Transaction) error {
	if tx.TransactionType == deposit {
		return e.repository.RegisterTransaction(ctx, tx)
	}
	tx.Amount = tx.Amount * -1
	return e.repository.RegisterTransaction(ctx, tx)

}
