package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/kernelquorum/city-counter/internal/config"
	internal_http "github.com/kernelquorum/city-counter/internal/http"
	"github.com/kernelquorum/city-counter/internal/repository"

	"net/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	var cityRepo repository.CityRepository = repository.NewGeoNames(cfg.GeoNamesURL)
	if cfg.CacheEnabled {
		cityRepo = repository.NewCityCacheWithRepository(cityRepo, cfg.CacheTTL)
	}

	handler := internal_http.Handler{
		CityRepository: cityRepo,
	}
	server := &http.Server{
		Addr:    ":" + strconv.Itoa(cfg.Port),
		Handler: internal_http.NewRouter(handler),
	}

	/* catch SIGTERM to gracefully shutdown server*/
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErr:
		if !errors.Is(err, http.ErrServerClosed) {
			log.Printf("http server: %v", err)
		}
		return
	case <-ctx.Done():
		stop()
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown: %v", err)
		server.Close()
	}
}
