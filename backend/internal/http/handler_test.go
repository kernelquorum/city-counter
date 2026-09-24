package http_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	internal_http "github.com/kernelquorum/city-counter/internal/http"
	"github.com/kernelquorum/city-counter/internal/model"
	"github.com/stretchr/testify/require"
)

type cityRepositoryStub struct {
	err error
}

func (r cityRepositoryStub) GetCities(context.Context) (model.CitiesResponse, error) {
	return model.CitiesResponse{Cities: []model.City{
		{Name: "London"}, {Name: "lisbon"}, {Name: "Łódź"},
	}}, r.err
}

func TestCountCities(t *testing.T) {
	tests := []struct {
		name   string
		letter string
		err    error
		status int
		body   string
	}{
		{name: "counts cities", letter: "L", status: http.StatusOK, body: `{"count":2}`},
		{name: "lowercase letter", letter: "l", status: http.StatusOK, body: `{"count":2}`},
		{name: "Unicode letter", letter: "ł", status: http.StatusOK, body: `{"count":1}`},
		{name: "no matches", letter: "Z", status: http.StatusOK, body: `{"count":0}`},
		{name: "missing letter", status: http.StatusBadRequest, body: "first_letter must be exactly one letter\n"},
		{name: "multiple letters", letter: "LL", status: http.StatusBadRequest, body: "first_letter must be exactly one letter\n"},
		{name: "invalid UTF-8", letter: "\xff", status: http.StatusBadRequest, body: "first_letter must be exactly one letter\n"},
		{name: "not a letter", letter: "1", status: http.StatusBadRequest, body: "first_letter must be a letter\n"},
		{name: "service error", letter: "L", err: errors.New("repository failed"), status: http.StatusInternalServerError, body: "failed to count cities\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := internal_http.Handler{
				CityRepository: cityRepositoryStub{err: tt.err},
			}
			target := "/cities/count"
			if tt.letter != "" {
				target += "?first_letter=" + url.QueryEscape(tt.letter)
			}
			req := httptest.NewRequest(http.MethodGet, target, nil)
			res := httptest.NewRecorder()

			h.CountCities(res, req)

			require.Equal(t, tt.status, res.Code)
			if tt.status == http.StatusOK {
				require.Equal(t, "application/json", res.Header().Get("Content-Type"))
				require.JSONEq(t, tt.body, res.Body.String())
			} else {
				require.Equal(t, tt.body, res.Body.String())
			}
		})
	}
}
