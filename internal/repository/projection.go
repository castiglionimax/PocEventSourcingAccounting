package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/castiglionimax/PocEventSourcingAccounting/internal/domain"
)

type (
	ProjectionAccount struct {
		mysql *sql.DB
	}
)

const (
	insertAccount       = "INSERT INTO accounts (id, name, amount, number, last_updated) VALUES (?,?, ?, ?,?)"
	UpdateAccountAmount = "UPDATE accounts SET amount = amount + ? WHERE id = ?;"
)

func NewProjection(db *sql.DB) *ProjectionAccount {
	return &ProjectionAccount{db}
}

func (p ProjectionAccount) CreateAccount(ctx context.Context, account domain.Account) error {
	insertStatement, err := p.mysql.Prepare(insertAccount)
	if err != nil {
		return err
	}

	_, err = insertStatement.Exec(account.ID, account.Name, 0, account.Number, time.Now().UTC())
	if err != nil {
		return err
	}
	return nil
}

func (p ProjectionAccount) RegisterTransaction(ctx context.Context, tx domain.Transaction) error {
	insertStatement, err := p.mysql.Prepare(UpdateAccountAmount)
	if err != nil {
		return err
	}

	_, err = insertStatement.Exec(tx.Amount, tx.AccountID)
	if err != nil {
		return err
	}
	return nil
}
