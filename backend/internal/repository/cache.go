package repository

import (
	"context"
	"sync"
	"time"

	"github.com/kernelquorum/city-counter/internal/model"
)

type CityCache struct {
	geoNamesRepo CityRepository
	validFor     time.Duration

	response  model.CitiesResponse
	expiresAt time.Time
	mtx       sync.Mutex
}

func NewCityCache(url string, ttl time.Duration) *CityCache {
	return NewCityCacheWithRepository(NewGeoNames(url), ttl)
}

// NewCityCacheWithRepository caches responses from source for the supplied TTL.
func NewCityCacheWithRepository(source CityRepository, ttl time.Duration) *CityCache {
	return &CityCache{
		geoNamesRepo: source,
		expiresAt:    time.Now(),
		validFor:     ttl,
	}
}

func (c *CityCache) GetCities(ctx context.Context) (model.CitiesResponse, error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	if time.Now().Before(c.expiresAt) {
		rsp := c.response
		return rsp, nil
	}

	rsp, err := c.geoNamesRepo.GetCities(ctx)
	if err != nil {
		return rsp, err
	}

	c.response = rsp
	c.expiresAt = time.Now().Add(c.validFor)

	return rsp, err
}
