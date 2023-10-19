package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"io"
	"net/http"

	"github.com/go-chi/render"

	"github.com/castiglionimax/PocEventSourcingAccounting/internal/domain"
	pkgError "github.com/castiglionimax/PocEventSourcingAccounting/pkg/error"
)

type (
	Service interface {
		CreateAccount(ctx context.Context, account domain.Account) (domain.AccountID, error)
		Transaction(ctx context.Context, transaction domain.Transaction) error
		GetBalance(ctx context.Context, accountID domain.AccountID) (float32, error)
	}

	Controller struct {
		service Service
	}
)

func NewController(service Service) (*Controller, error) {
	if service == nil {
		return nil, errors.New("service should not be nil")
	}
	return &Controller{
		service: service,
	}, nil
}

func (c Controller) CreateAccount(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, pkgError.ErrReadingBody.Error(), http.StatusBadRequest)
		return
	}

	var req domain.Account

	if err = json.Unmarshal(data, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	account, err := c.service.CreateAccount(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, account)
	w.WriteHeader(http.StatusCreated)
}

func (c Controller) Transaction(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, pkgError.ErrReadingBody.Error(), http.StatusBadRequest)
		return
	}

	var req domain.Transaction

	if err = json.Unmarshal(data, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = c.service.Transaction(r.Context(), req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c Controller) Balance(w http.ResponseWriter, r *http.Request) {
	accountID := chi.URLParam(r, "id")
	amount, err := c.service.GetBalance(r.Context(), domain.AccountID(accountID))
	if err != nil {
		if (errors.As(err, &pkgError.HandlerError{})) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write([]byte(fmt.Sprintf("amount:%f", amount)))

	w.WriteHeader(http.StatusOK)
}
