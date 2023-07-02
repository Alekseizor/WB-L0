package orders

import (
	"context"
	"database/sql"
	"sync"
)

type RepoOrderInMemory struct {
	orders map[string]*OrderAllData
	mu     *sync.RWMutex
}

func NewRepoOrderInMemory() (*RepoOrderInMemory, error) {
	return &RepoOrderInMemory{
		orders: make(map[string]*OrderAllData, 0),
		mu:     &sync.RWMutex{},
	}, nil
}

func (om *RepoOrderInMemory) AddOrder(item OrderAllData, ctx context.Context) error {
	om.mu.Lock()
	om.orders[item.OrderUID] = &item
	om.mu.Unlock()
	return nil
}

func (om *RepoOrderInMemory) GetOrderByID(orderUID string, ctx context.Context) (*OrderAllData, error) {
	om.mu.RLock()
	link, existence := om.orders[orderUID]
	om.mu.RUnlock()
	if existence {
		return link, nil
	}
	return nil, sql.ErrNoRows
}

func (om *RepoOrderInMemory) GetAll(ctx context.Context) ([]*OrderAllData, error) {
	om.mu.RLock()
	orders := make([]*OrderAllData, 0, len(om.orders))
	index := 0
	for _, order := range om.orders {
		orders[index] = order
		index++
	}
	om.mu.RUnlock()
	return orders, nil
}
