package db

import "time"

type Transaction struct {
	ID        string    `json:"id"`
	Kind      string    `json:"kind"`
	Status    string    `json:"status"`
	Account   string    `json:"account"`
	AssetCode string    `json:"asset_code"`
	Amount    string    `json:"amount,omitempty"`
	AmountIn  string    `json:"amount_in,omitempty"`
	AmountOut string    `json:"amount_out,omitempty"`
	AmountFee string    `json:"amount_fee,omitempty"`
	URL       string    `json:"url,omitempty"`
	StartedAt time.Time `json:"started_at"`
	UpdatedAt time.Time `json:"updated_at"`
	KYCFields []string  `json:"kyc_fields,omitempty"`
}

type TransactionStore interface {
	Create(tx Transaction) error
	GetByID(id string) (Transaction, bool)
	ListByAccount(account string, limit int, cursor string) []Transaction
	Update(tx Transaction) error
	UpdateStatus(id string, status string, updatedAt time.Time) (Transaction, bool)
}

type Customer struct {
	Account string            `json:"account"`
	Fields  map[string]string `json:"fields"`
}

type CustomerStore interface {
	Put(customer Customer) error
	Get(account string) (Customer, bool)
}
