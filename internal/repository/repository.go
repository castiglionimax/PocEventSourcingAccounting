package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	pkgError "github.com/castiglionimax/challengeXepelin/pkg/error"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/castiglionimax/challengeXepelin/internal/domain"
)

type (
	producer interface {
		Produce(msg *kafka.Message, deliveryChan chan kafka.Event) error
	}

	Repository struct {
		producer producer
		mongo    *mongo.Client
		topic    string
		mysql    *sql.DB
	}
)

func NewRepository(producer producer, topic string, mongo *mongo.Client, mysql *sql.DB) *Repository {
	return &Repository{producer: producer, topic: topic, mongo: mongo, mysql: mysql}
}

const (
	createAccount  = "account_created"
	saveWithdrawal = "withdrawal_saved"
	saveDeposit    = "deposit_saved"

	getAmount = "SELECT amount FROM accounts WHERE id = ?"
)

func (r Repository) CreateAccount(ctx context.Context, account domain.Account) (domain.AccountID, error) {
	account.ID = domain.AccountID(uuid.New().String())
	eventModel := NewModel(createAccount, account.ID.String(), account)

	if err := r.apply(ctx, eventModel); err != nil {
		return "", err
	}
	return account.ID, nil
}

func (r Repository) SaveTransaction(ctx context.Context, transaction domain.Transaction) error {
	var eventModel Model
	if transaction.TransactionType == "deposit" {
		eventModel = NewModel(saveDeposit, transaction.AccountID.String(), transaction)
	} else {
		eventModel = NewModel(saveWithdrawal, transaction.AccountID.String(), transaction)
	}
	return r.apply(ctx, eventModel)
}

func (r Repository) apply(ctx context.Context, event Model) error {
	coll := r.mongo.Database("event_store").Collection("accounts")

	session, err := r.mongo.StartSession()
	if err != nil {
		return err
	}

	defer session.EndSession(ctx)

	if err = session.StartTransaction(); err != nil {
		return err
	}

	return mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		_, err = coll.InsertOne(ctx, event)
		if err != nil {
			return err
		}

		enAccount, err := json.Marshal(event)
		if err != nil {
			return err
		}

		deliveryChan := make(chan kafka.Event, 10000)
		err = r.producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &r.topic, Partition: kafka.PartitionAny},
			Value:          enAccount},
			deliveryChan,
		)

		if err != nil {
			_ = session.AbortTransaction(ctx)
			return err
		}

		if err = session.CommitTransaction(sc); err != nil {
			return err
		}
		return nil
	})
}

func (r Repository) GetBalance(ctx context.Context, accountID domain.AccountID) (float32, error) {
	row := r.mysql.QueryRow(getAmount, accountID)
	var amount float32
	err := row.Scan(&amount)
	if err != nil {
		return 0, pkgError.HandlerError{Cause: errors.New("not found")}
	}
	return amount, nil
}
