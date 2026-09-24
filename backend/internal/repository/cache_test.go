package repository_test

import (
	"context"
	"testing"
	"testing/synctest"
	"time"

	"github.com/kernelquorum/city-counter/internal/model"
	"github.com/kernelquorum/city-counter/internal/repository"
	"github.com/stretchr/testify/require"
)

type cityRepositoryFunc func(context.Context) (model.CitiesResponse, error)

func (f cityRepositoryFunc) GetCities(ctx context.Context) (model.CitiesResponse, error) {
	return f(ctx)
}

func TestCityCacheExpiry(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		calls := 0
		city := "Shanghai"
		source := cityRepositoryFunc(func(context.Context) (model.CitiesResponse, error) {
			calls++
			return model.CitiesResponse{Cities: []model.City{{Name: city}}}, nil
		})
		ttl := 10 * time.Second
		cache := repository.NewCityCacheWithRepository(source, ttl)

		first, err := cache.GetCities(t.Context())
		require.NoError(t, err)
		require.Equal(t, "Shanghai", first.Cities[0].Name)
		require.Equal(t, 1, calls)

		// The source changes, but the cached response is still fresh.
		city = "Beijing"
		time.Sleep(ttl - time.Nanosecond)
		cached, err := cache.GetCities(t.Context())
		require.NoError(t, err)
		require.Equal(t, first, cached)
		require.Equal(t, 1, calls)

		// At exactly the configured TTL, the cache must fetch the new response.
		time.Sleep(time.Nanosecond)
		refreshed, err := cache.GetCities(t.Context())
		require.NoError(t, err)
		require.Len(t, refreshed.Cities, 1)
		require.Equal(t, "Beijing", refreshed.Cities[0].Name)
		require.Equal(t, 2, calls)

		cached, err = cache.GetCities(t.Context())
		require.NoError(t, err)
		require.Equal(t, refreshed, cached)
		require.Equal(t, 2, calls)
	})
}
