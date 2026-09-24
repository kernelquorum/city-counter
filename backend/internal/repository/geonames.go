package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/kernelquorum/city-counter/internal/model"
)

type GeoNames struct {
	url    string
	client *http.Client
}

func NewGeoNames(url string) *GeoNames {
	return NewGeoNamesWithClient(url, &http.Client{Timeout: 60 * time.Second})
}

// NewGeoNamesWithClient creates a repository using the supplied HTTP client.
// If client is nil, it uses the default client configuration.
func NewGeoNamesWithClient(url string, client *http.Client) *GeoNames {
	if client == nil {
		client = &http.Client{Timeout: 60 * time.Second}
	}
	return &GeoNames{url: url, client: client}
}

func (g *GeoNames) GetCities(ctx context.Context) (model.CitiesResponse, error) {
	response := model.CitiesResponse{}

	err := g.exponentalBackoff(ctx, 5, func() (error, bool) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, g.url, nil)
		if err != nil {
			return err, false
		}

		resp, err := g.client.Do(req)
		if err != nil {
			return err, true
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("http get returned error: %d", resp.StatusCode), g.shouldRetry(resp.StatusCode)
		}

		err = json.NewDecoder(resp.Body).Decode(&response)
		return err, false
	})

	return response, err
}

func (g *GeoNames) shouldRetry(status int) bool {
	switch status {
	case
		http.StatusRequestTimeout,      // 408
		http.StatusTooManyRequests,     // 429
		http.StatusInternalServerError, // 500
		http.StatusBadGateway,          // 502
		http.StatusServiceUnavailable,  // 503
		http.StatusGatewayTimeout:      // 504
		return true
	default:
		return false
	}
}

func (g *GeoNames) exponentalBackoff(ctx context.Context, maxRetries int, fn func() (error, bool)) error {
	delay := 500 * time.Millisecond
	for attempt := range maxRetries {
		err, retry := fn()
		if !retry {
			return err
		}

		if attempt == maxRetries-1 {
			break
		}

		timer := time.NewTimer(delay)
		defer timer.Stop()

		select {
		case <-timer.C:
			delay *= 2
			continue
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return fmt.Errorf("max number of retires reached")
}
