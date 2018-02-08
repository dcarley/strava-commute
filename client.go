package main

import (
	"net/http"
	"time"

	strava "github.com/strava/go.strava"
)

// Transport can be overridden for the purpose of testing.
var Transport http.RoundTripper = &http.Transport{}

func NewClient(token string) *strava.Client {
	httpClient := &http.Client{
		Transport: Transport,
		Timeout:   2 * time.Second,
	}

	return strava.NewClient(token, httpClient)
}
