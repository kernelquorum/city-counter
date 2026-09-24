package service_test

import (
	"context"
	"testing"

	"github.com/kernelquorum/city-counter/internal/model"
	"github.com/kernelquorum/city-counter/internal/service"
	"github.com/stretchr/testify/require"
)

type cityRepositoryStub struct {
	cities model.CitiesResponse
}

func (r cityRepositoryStub) GetCities(context.Context) (model.CitiesResponse, error) {
	return r.cities, nil
}

func TestCountByStartingLetter(t *testing.T) {
	tests := []struct {
		name   string
		cities []string
		letter rune
		want   int
	}{
		{
			name:   "counts matching cities",
			cities: []string{"Shanghai", "Sydney", "Beijing", "Oslo"},
			letter: 'S', want: 2,
		},
		{
			name:   "compares letters regardless of case",
			cities: []string{"London", "lisbon", "Berlin"},
			letter: 'l', want: 2,
		},
		{
			name:   "decodes Unicode and compares case",
			cities: []string{"Łódź", "łowicz", "London", "Évora"},
			letter: 'Ł', want: 2,
		},
		{
			name:   "decodes non Latin first rune",
			cities: []string{"東京", "東大阪", "北京"},
			letter: '東', want: 2,
		},
		{
			name:   "accented letters are distinct",
			cities: []string{"Évora", "évian", "Edinburgh"},
			letter: 'é', want: 2,
		},
		{
			name:   "ignores empty names and later letters",
			cities: []string{"", "Oslo", "Berlin"},
			letter: 'l', want: 0,
		},
		{
			name:   "empty city list",
			letter: 'A', want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := model.CitiesResponse{}
			for _, name := range tt.cities {
				response.Cities = append(response.Cities, model.City{Name: name})
			}
			cs := service.CityService{
				CityRepository: cityRepositoryStub{cities: response},
			}

			count, err := cs.CountByStartingLetter(t.Context(), tt.letter)

			require.NoError(t, err)
			require.Equal(t, tt.want, count)
		})
	}
}
