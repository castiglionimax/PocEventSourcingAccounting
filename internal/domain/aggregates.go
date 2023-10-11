package domain

type Transaction struct {
	AccountID       AccountID `json:"account_id"`
	TransactionType string    `json:"transaction_type"`
	Amount          float32   `json:"amount"`
}
