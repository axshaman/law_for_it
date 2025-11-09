package payment

import (
	"context"
	"fmt"
	"time"

	"law_for_it/internal/errs"
	"law_for_it/internal/model"
	"law_for_it/internal/state"
)

type Service struct {
	manager *state.Manager
}

func NewService(manager *state.Manager, currency string) (*Service, error) {
	if manager == nil {
		return nil, fmt.Errorf("state manager is required")
	}
	if err := manager.UpdateAccountCurrency(currency); err != nil {
		return nil, err
	}
	return &Service{manager: manager}, nil
}

func (s *Service) TopUp(ctx context.Context, amount float64, message string) (model.Account, error) {
	if amount <= 0 {
		return model.Account{}, errs.ErrInvalidAmount
	}
	var account model.Account
	err := s.manager.Update(func(data *state.Data) error {
		data.Account.Balance += amount
		data.Account.Updated = time.Now().UTC()
		txn := model.Transaction{
			ID:        data.NextTransactionID,
			Type:      "topup",
			Amount:    amount,
			Balance:   data.Account.Balance,
			Message:   message,
			CreatedAt: time.Now().UTC(),
		}
		data.NextTransactionID++
		prependTransaction(&data.Transactions, txn)
		account = data.Account
		return nil
	})
	if err != nil {
		return model.Account{}, err
	}
	return account, nil
}

func (s *Service) Charge(ctx context.Context, amount float64, message string) error {
	if amount <= 0 {
		return errs.ErrInvalidAmount
	}
	return s.manager.Update(func(data *state.Data) error {
		if data.Account.Balance < amount {
			return errs.ErrInsufficientBalance
		}
		data.Account.Balance -= amount
		data.Account.Updated = time.Now().UTC()
		txn := model.Transaction{
			ID:        data.NextTransactionID,
			Type:      "charge",
			Amount:    -amount,
			Balance:   data.Account.Balance,
			Message:   message,
			CreatedAt: time.Now().UTC(),
		}
		data.NextTransactionID++
		prependTransaction(&data.Transactions, txn)
		return nil
	})
}

func (s *Service) Balance(ctx context.Context) (model.Account, error) {
	var account model.Account
	err := s.manager.View(func(data *state.Data) error {
		account = data.Account
		return nil
	})
	return account, err
}

func (s *Service) History(ctx context.Context, limit int) ([]model.Transaction, error) {
	if limit <= 0 {
		limit = 50
	}
	transactions := []model.Transaction{}
	err := s.manager.View(func(data *state.Data) error {
		if limit > len(data.Transactions) {
			limit = len(data.Transactions)
		}
		transactions = append(transactions, data.Transactions[:limit]...)
		return nil
	})
	return transactions, err
}

func (s *Service) Currency() string {
	return s.manager.Account().Currency
}

func prependTransaction(target *[]model.Transaction, txn model.Transaction) {
	const maxTransactions = 200
	*target = append([]model.Transaction{txn}, *target...)
	if len(*target) > maxTransactions {
		*target = (*target)[:maxTransactions]
	}
}
