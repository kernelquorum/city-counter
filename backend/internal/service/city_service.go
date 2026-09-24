package service

import (
	"context"
	"unicode"
	"unicode/utf8"

	"github.com/kernelquorum/city-counter/internal/repository"
)

type CityService struct {
	repository.CityRepository
}

func (cs *CityService) CountByStartingLetter(ctx context.Context, letter rune) (int, error) {
	cities, err := cs.CityRepository.GetCities(ctx)
	if err != nil {
		return 0, err
	}

	found := 0

	for _, city := range cities.Cities {
		first, size := utf8.DecodeRuneInString(city.Name)
		if size == 0 {
			continue
		}

		if unicode.ToLower(first) == unicode.ToLower(letter) {
			found++
		}
	}

	return found, nil
}
