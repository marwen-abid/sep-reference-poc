package db

import (
	"sort"
	"sync"
	"time"
)

type MemoryTransactionStore struct {
	mu  sync.RWMutex
	txs map[string]Transaction
}

type MemoryCustomerStore struct {
	mu        sync.RWMutex
	customers map[string]Customer
}

func NewMemoryTransactionStore() *MemoryTransactionStore {
	return &MemoryTransactionStore{txs: map[string]Transaction{}}
}

func NewMemoryCustomerStore() *MemoryCustomerStore {
	return &MemoryCustomerStore{customers: map[string]Customer{}}
}

func (s *MemoryTransactionStore) Create(tx Transaction) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.txs[tx.ID] = tx
	return nil
}

func (s *MemoryTransactionStore) GetByID(id string) (Transaction, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	tx, ok := s.txs[id]
	return tx, ok
}

func (s *MemoryTransactionStore) ListByAccount(account string, limit int, cursor string) []Transaction {
	s.mu.RLock()
	defer s.mu.RUnlock()

	items := make([]Transaction, 0, len(s.txs))
	for _, tx := range s.txs {
		if tx.Account == account {
			items = append(items, tx)
		}
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].StartedAt.After(items[j].StartedAt)
	})

	if limit <= 0 || limit > len(items) {
		limit = len(items)
	}
	return items[:limit]
}

func (s *MemoryTransactionStore) Update(tx Transaction) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.txs[tx.ID] = tx
	return nil
}

func (s *MemoryTransactionStore) UpdateStatus(id string, status string, updatedAt time.Time) (Transaction, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	tx, ok := s.txs[id]
	if !ok {
		return Transaction{}, false
	}
	tx.Status = status
	tx.UpdatedAt = updatedAt
	s.txs[id] = tx
	return tx, true
}

func (s *MemoryCustomerStore) Put(customer Customer) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.customers[customer.Account] = customer
	return nil
}

func (s *MemoryCustomerStore) Get(account string) (Customer, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.customers[account]
	return c, ok
}
