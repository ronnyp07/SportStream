package cchttp

import (
	"io"
	"net/http"
	"time"
)

type Client interface {
	Get(url string) (*http.Response, error)
	Post(url, contentType string, body io.Reader) (*http.Response, error)
	Do(req *http.Request) (*http.Response, error)
}

func NewClient(maxIdleConns, maxConnsPerHost, maxIdleConnsPerHost int, timeout time.Duration) Client {
	transport := &http.Transport{
		MaxIdleConns:        maxIdleConns,
		MaxConnsPerHost:     maxConnsPerHost,
		MaxIdleConnsPerHost: maxIdleConnsPerHost,
	}

	return &http.Client{
		Transport: transport,
		Timeout:   timeout,
	}
}
