package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config represents service configuration for dp-maps-api
type Config struct {
	BindAddr                   string        `envconfig:"BIND_ADDR"`
	GracefulShutdownTimeout    time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	HealthCheckInterval        time.Duration `envconfig:"HEALTHCHECK_INTERVAL"`
	HealthCheckCriticalTimeout time.Duration `envconfig:"HEALTHCHECK_CRITICAL_TIMEOUT"`
	MapsAPIURL                 string        `envconfig:"MAPS_API_URL"`
	OrdnanceSurveyAPIURL       string        `envconfig:"ORDNANCE_SURVEY_API_URL"`
	OrdnanceSurveyAPIKey       string        `envconfig:"ORDNANCE_SURVEY_API_KEY"           json:"-"`
	CacheMaxAge                time.Duration `envconfig:"CACHE_MAX_AGE"`
}

var cfg *Config

// Get returns the default config with any modifications through environment
// variables
func Get() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg = &Config{
		BindAddr:                   "localhost:27900",
		GracefulShutdownTimeout:    5 * time.Second,
		HealthCheckInterval:        30 * time.Second,
		HealthCheckCriticalTimeout: 90 * time.Second,
		MapsAPIURL:                 "http://localhost:27900/",
		OrdnanceSurveyAPIURL:       "https://api.os.uk/",
		CacheMaxAge:                24 * time.Hour,
	}

	return cfg, envconfig.Process("", cfg)
}
