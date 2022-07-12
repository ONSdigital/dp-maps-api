package api

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/ONSdigital/log.go/v2/log"
)

const (
	osApiKeyHeader = "key"
)

// CreateOSMapsProxy returns a ReverseProxy that modifies requests to add an Ordnance Survey api key and forward them on
// to the O/S api. It can also optionally modify the response (for example the URLS contained in the response)
func CreateOSMapsProxy(target *url.URL, cfg Config, responseModifier func(r *http.Response) error) *httputil.ReverseProxy {
	targetQuery := target.RawQuery
	hasResponseModifier := responseModifier != nil
	director := func(req *http.Request) {
		log.Info(req.Context(), "proxying request", log.HTTP(req, 0, 0, nil, nil), log.Data{
			"target": target,
		})
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.Host = req.URL.Host
		//req.RequestURI = req.URL.Path
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
		req.Header.Set(osApiKeyHeader, cfg.OrdnanceSurveyAPIKey) // Add O/S Api Key
		if hasResponseModifier {
			req.Header.Del("Accept-Encoding")
		} // Clear the encoding so responses aren't gzipped and can be modified by a response modifier more easily
	}

	return &httputil.ReverseProxy{
		Director:       director,
		ModifyResponse: CacheModifier(cfg, responseModifier),
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       180 * time.Second,
			TLSHandshakeTimeout:   5 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
}

// CacheModifier returns a ResponseModifier that sets a maximum cache header
// It optionally takes another response modifier that it wraps around and performs first when called
func CacheModifier(cfg Config, parentModifier func(r *http.Response) error) func(*http.Response) error {
	maxAge := fmt.Sprintf("%.f", cfg.CacheMaxAge.Seconds())
	return func(response *http.Response) error {
		if parentModifier != nil {
			err := parentModifier(response)
			if err != nil {
				return err
			}
		}
		if response.Header == nil {
			response.Header = http.Header{}
		}
		response.Header.Set("Cache-Control", "max-age="+maxAge)
		return nil
	}
}

// StringReplaceResponseModifier returns a function that takes in a pointer to a Response object that modifies that
// response by changing every instance of a string in the body to another string.
func StringReplaceResponseModifier(from, to string) func(*http.Response) error {
	return func(r *http.Response) error {
		ctx := r.Request.Context()
		log.Info(ctx, "rewriting response", log.Data{
			"url":  r.Request.URL.String(),
			"from": from,
			"to":   to,
		})
		origBodyReader := r.Body

		body, err := io.ReadAll(origBodyReader)
		if err != nil {
			return err
		}
		origBodyReader.Close()

		body = bytes.Replace(body, []byte(from), []byte(to), -1)

		r.Body = io.NopCloser(bytes.NewReader(body))
		return nil
	}
}
