package main

import (
	"net/http"
	"strconv"
	"time"
)

// RateLimitedTransport is used to conform to rate limits that are
// communicated through "X-RateLimit-" headers, like GitHub's API. It
// implements http.RoundTripper and can be used for configuring a http.Client.
type RateLimitedTransport struct {
	Base http.RoundTripper
}

func (t *RateLimitedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	res, err := t.base().RoundTrip(req)
	if err != nil {
		return res, err
	}

	// Fetch headers
	remStr := res.Header.Get("X-RateLimit-Remaining")
	if remStr == "" {
		return res, err
	}
	resetStr := res.Header.Get("X-RateLimit-Reset")
	if resetStr == "" {
		return res, err
	}

	rem, err := strconv.Atoi(remStr)
	if err != nil {
		return res, err
	}
	epoch, err := strconv.ParseInt(resetStr, 10, 64)
	if err != nil {
		return res, err
	}
	reset := time.Unix(epoch, 0)

	// Determine sleep time
	untilReset := reset.Sub(time.Now())
	delay := time.Duration(float64(untilReset) / (float64(rem) + 1))
	time.Sleep(delay)

	return res, err
}

// base returns the wrapped RoundTripper, using the default if no custom
// RoundTripper has been provided.
func (t *RateLimitedTransport) base() http.RoundTripper {
	if t.Base != nil {
		return t.Base
	}
	return http.DefaultTransport
}
