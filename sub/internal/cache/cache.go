package cache

import (
	"context"
	"errors"
	"sync"

	"test/internal/model"
	pg "test/internal/pgconn"
)

type Cache struct {
	cache map[string]model.OrderMessage
	mutex *sync.Mutex
}

func InitCache() *Cache {

	var mtx sync.Mutex

	return &Cache{
		make(map[string]model.OrderMessage),
		&mtx,
	}
}

func (c *Cache) Add(order model.OrderMessage) {
	c.mutex.Lock()
	c.cache[order.OrderUID] = order
	c.mutex.Unlock()
}

func (c *Cache) Get(id string) (model.OrderMessage, error) {

	defer c.mutex.Unlock()
	c.mutex.Lock()
	order, ok := c.cache[id]

	if !ok {
		return model.OrderMessage{}, errors.New("order not found")
	}

	return order, nil
}

func (c *Cache) CacheFromDb(ctx context.Context, db *pg.DB) error {

	orders, err := db.GetAllMessages()

	if err != nil {
		return err
	}

	for _, item := range orders {
		c.Add(item)
	}

	return nil
}
