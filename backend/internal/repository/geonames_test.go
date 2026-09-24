package repository_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"testing/synctest"
	"time"

	"github.com/kernelquorum/city-counter/internal/repository"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type GeoNamesSuite struct {
	suite.Suite

	mux    *http.ServeMux
	server *httptest.Server
	client *repository.GeoNames
}

func (s *GeoNamesSuite) SetupTest() {
	s.mux = http.NewServeMux()
	s.server = httptest.NewServer(s.mux)
	s.client = repository.NewGeoNames(s.server.URL)
}

func (s *GeoNamesSuite) TearDownTest() {
	s.server.Close()
}

func TestGeoNames(t *testing.T) {
	suite.Run(t, new(GeoNamesSuite))
}

func (s *GeoNamesSuite) TestGetCitiesSingle() {
	body, err := os.ReadFile("testdata/geonames_single.json")
	s.Require().NoError(err)

	s.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		s.Equal(http.MethodGet, r.Method)
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write(body)
		s.NoError(err)
	})

	response, err := s.client.GetCities(s.T().Context())

	s.Require().NoError(err)
	s.Require().Len(response.Cities, 1)
	s.Equal("Shanghai", response.Cities[0].Name)
}

func (s *GeoNamesSuite) TestGetCitiesTwo() {
	body, err := os.ReadFile("testdata/geonames_two.json")
	s.Require().NoError(err)

	s.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		s.Equal(http.MethodGet, r.Method)
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write(body)
		s.NoError(err)
	})

	response, err := s.client.GetCities(s.T().Context())

	s.Require().NoError(err)
	s.Require().Len(response.Cities, 2)
	s.Equal("Shanghai", response.Cities[0].Name)
	s.Equal("Beijing", response.Cities[1].Name)
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func (s *GeoNamesSuite) TestGetCitiesRetries() {
	body, err := os.ReadFile("testdata/geonames_single.json")
	s.Require().NoError(err)

	synctest.Test(s.T(), func(t *testing.T) {
		start := time.Now()
		var attempts []time.Duration
		client := repository.NewGeoNamesWithClient("http://geonames.test", &http.Client{
			Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				attempts = append(attempts, time.Since(start))
				status := http.StatusServiceUnavailable
				if len(attempts) == 3 {
					status = http.StatusOK
				}
				return &http.Response{
					StatusCode: status,
					Header:     make(http.Header),
					Body:       io.NopCloser(bytes.NewReader(body)),
					Request:    r,
				}, nil
			}),
		})

		response, err := client.GetCities(t.Context())

		require.NoError(t, err)
		require.Equal(t, []time.Duration{0, 500 * time.Millisecond, 1500 * time.Millisecond}, attempts)
		require.Equal(t, 1500*time.Millisecond, time.Since(start))
		require.Len(t, response.Cities, 1)
		require.Equal(t, "Shanghai", response.Cities[0].Name)
	})
}

func (s *GeoNamesSuite) TestGetCitiesDoesNotRetry() {
	synctest.Test(s.T(), func(t *testing.T) {
		start := time.Now()
		requests := 0
		client := repository.NewGeoNamesWithClient("http://geonames.test", &http.Client{
			Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				requests++
				return &http.Response{
					StatusCode: http.StatusBadRequest,
					Header:     make(http.Header),
					Body:       http.NoBody,
					Request:    r,
				}, nil
			}),
		})

		_, err := client.GetCities(t.Context())

		require.Error(t, err)
		require.Equal(t, 1, requests)
		require.Zero(t, time.Since(start))
	})
}
