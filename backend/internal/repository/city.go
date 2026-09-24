package repository

import (
	"context"

	"github.com/kernelquorum/city-counter/internal/model"
)

type CityRepository interface {
	GetCities(ctx context.Context) (model.CitiesResponse, error)
}
