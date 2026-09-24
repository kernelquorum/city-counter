package http

import (
	"encoding/json"
	"log"
	"net/http"
	"unicode"
	"unicode/utf8"

	"github.com/kernelquorum/city-counter/internal/service"
)

type Handler struct {
	service.CityService
}

type CountCitiesResponse struct {
	Count int `json:"count"`
}

func (h *Handler) CountCities(w http.ResponseWriter, r *http.Request) {
	queryArg := r.URL.Query().Get("first_letter")
	letters := []rune(queryArg)

	if !utf8.ValidString(queryArg) || len(letters) != 1 {
		http.Error(w, "first_letter must be exactly one letter", http.StatusBadRequest)
		return
	}

	if !unicode.IsLetter(letters[0]) {
		http.Error(w, "first_letter must be a letter", http.StatusBadRequest)
		return
	}

	count, err := h.CityService.CountByStartingLetter(r.Context(), letters[0])
	if err != nil {
		http.Error(w, "failed to count cities", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(CountCitiesResponse{Count: count}); err != nil {
		log.Printf("citou count encode error: %v", err)
		return
	}
}
