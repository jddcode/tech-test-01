package cache

import (
	"errors"
	"github.com/jddcode/tech-test-ennismore/internal/handler-weather/structs"
	"sync"
)

//go:generate mockgen -destination=../../mocks/mock-cache.go -package=mocks . Cache
type Cache interface {
	Get(city string) ([]structs.ResultForecast, error)
	Store(city string, prediction []structs.ResultForecast)
}

type cache struct {
	content map[string][]structs.ResultForecast
	lock    sync.RWMutex
}

func (c *cache) Store(city string, predictions []structs.ResultForecast) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.content[city] = predictions
}

func (c cache) Get(city string) ([]structs.ResultForecast, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	val, exists := c.content[city]
	if !exists {
		return []structs.ResultForecast{}, errors.New("cache miss")
	}
	return val, nil
}
