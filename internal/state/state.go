package state

import "law_for_it/internal/model"

type Data struct {
	NextDocumentID    uint                `json:"next_document_id"`
	Documents         []model.Document    `json:"documents"`
	Account           model.Account       `json:"account"`
	NextTransactionID uint                `json:"next_transaction_id"`
	Transactions      []model.Transaction `json:"transactions"`
}

func New() *Data {
	return &Data{
		NextDocumentID:    1,
		Documents:         []model.Document{},
		Account:           model.Account{Balance: 0, Currency: "EUR"},
		NextTransactionID: 1,
		Transactions:      []model.Transaction{},
	}
}
