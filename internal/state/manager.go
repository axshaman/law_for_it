package state

import (
	"fmt"
	"sync"

	"law_for_it/internal/model"
	"law_for_it/internal/storage"
)

type Manager struct {
	store *storage.Store
	mu    sync.RWMutex
	data  *Data
}

func NewManager(store *storage.Store, defaults *Data) (*Manager, error) {
	if defaults == nil {
		defaults = New()
	}
	data := *defaults
	if err := store.Load(&data); err != nil {
		return nil, fmt.Errorf("failed to load state: %w", err)
	}
	normalise(&data)
	return &Manager{store: store, data: &data}, nil
}

func normalise(data *Data) {
	if data.NextDocumentID == 0 {
		data.NextDocumentID = 1
	}
	if data.NextTransactionID == 0 {
		data.NextTransactionID = 1
	}
	if data.Account.Currency == "" {
		data.Account.Currency = "EUR"
	}
}

func (m *Manager) View(fn func(*Data) error) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return fn(m.data)
}

func (m *Manager) Update(fn func(*Data) error) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if err := fn(m.data); err != nil {
		return err
	}
	return m.store.Save(m.data)
}

func (m *Manager) UpdateAccountCurrency(currency string) error {
	if currency == "" {
		return nil
	}
	return m.Update(func(data *Data) error {
		data.Account.Currency = currency
		return nil
	})
}

func (m *Manager) Account() model.Account {
	var account model.Account
	m.View(func(data *Data) error {
		account = data.Account
		return nil
	})
	return account
}
