package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"
)

const defaultGeoNamesURL = "http://api.geonames.org/searchJSON?featureClass=P&maxRows=50&orderby=population&username=hsample"

type Config struct {
	GeoNamesURL  string
	CacheEnabled bool
	CacheTTL     time.Duration
	Port         int
}

func Load() (*Config, error) {
	var cfg Config

	cfg.GeoNamesURL = env("GEONAMES_URL", defaultGeoNamesURL)
	u, err := url.Parse(cfg.GeoNamesURL)
	if err != nil || u.Hostname() == "" || (u.Scheme != "http" && u.Scheme != "https") {
		return &cfg, fmt.Errorf("GEONAMES_URL must be an absolute HTTP or HTTPS URL")
	}

	cfg.CacheEnabled, err = strconv.ParseBool(env("CACHE_ENABLED", "true"))
	if err != nil {
		return &cfg, fmt.Errorf("CACHE_ENABLED must be a boolean: %w", err)
	}

	cfg.CacheTTL, err = time.ParseDuration(env("CACHE_TTL", "1m"))
	if err != nil || cfg.CacheTTL <= 0 {
		return &cfg, fmt.Errorf("CACHE_TTL must be a positive duration, for example 30s or 5m")
	}

	cfg.Port, err = strconv.Atoi(env("PORT", "8080"))
	if err != nil || cfg.Port < 1 || cfg.Port > 65535 {
		return &cfg, fmt.Errorf("PORT must be an integer between 1 and 65535")
	}

	return &cfg, nil
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
