package http_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	internalhttp "github.com/kernelquorum/city-counter/internal/http"
	"github.com/stretchr/testify/require"
)

func TestRouter(t *testing.T) {
	router := internalhttp.NewRouter(internalhttp.Handler{})

	api := httptest.NewRecorder()
	router.ServeHTTP(api, httptest.NewRequest(http.MethodGet, "/cities/count", nil))
	require.Equal(t, http.StatusBadRequest, api.Code)
	require.Contains(t, api.Body.String(), "first_letter")

	page := httptest.NewRecorder()
	router.ServeHTTP(page, httptest.NewRequest(http.MethodGet, "/", nil))
	require.Equal(t, http.StatusNotFound, page.Code)
}
