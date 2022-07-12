package api

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"net/url"
	"time"
)

//API provides a struct to wrap the api around
type API struct {
	Router *mux.Router
}

// Config is a struct that contains configuration values required by the API
type Config struct {
	MapsAPIURL           string
	OrdnanceSurveyAPIURL *url.URL
	OrdnanceSurveyAPIKey string
	CacheMaxAge          time.Duration
}

//Setup function sets up the api and returns an api
func Setup(ctx context.Context, cfg Config, r *mux.Router) (*API, error) {
	api := &API{
		Router: r,
	}

	if cfg.OrdnanceSurveyAPIURL == nil {
		return nil, errors.New("nil OrdnanceSurveyAPIURL supplied in api setup config")
	}
	responseModifier := StringReplaceResponseModifier(cfg.OrdnanceSurveyAPIURL.String(), cfg.MapsAPIURL)
	modifyingProxy := CreateOSMapsProxy(cfg.OrdnanceSurveyAPIURL, cfg, responseModifier)
	nonModifyingProxy := CreateOSMapsProxy(cfg.OrdnanceSurveyAPIURL, cfg, nil)

	r.Handle("/maps/vector/v1/vts/resources/styles", modifyingProxy)
	r.Handle("/maps/vector/v1/vts", modifyingProxy)
	// Proxy Fallback
	r.Handle("/maps/vector/v1/vts/{uri:.*}", nonModifyingProxy)

	return api, nil
}
