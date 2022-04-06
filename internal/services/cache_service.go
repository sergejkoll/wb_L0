package services

import (
	"errors"
	"sync"

	"L0/internal/models"
)

type Cache struct {
	sync.RWMutex
	data map[string]models.Order
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string]models.Order),
	}
}

func (cache *Cache) Set(order models.Order) (err error) {
	cache.Lock()
	defer cache.Unlock()

	if _, ok := cache.data[order.OrderUID]; ok {
		return errors.New("id already exists")
	}
	cache.data[order.OrderUID] = order
	return nil
}

func (cache *Cache) Get(key string) (order models.Order, err error) {
	cache.RLock()
	defer cache.RUnlock()

	order, ok := cache.data[key]
	if !ok {
		return models.Order{}, errors.New("no such id exists")
	}
	return order, nil
}

func (cache *Cache) DataRecovery(db *DB) (err error) {
	cache.Lock()
	defer cache.Unlock()

	orders, err := db.GetOrders()
	if err != nil {
		return err
	}

	for _, order := range orders {
		cache.data[order.OrderUID] = order
	}

	return nil
}
