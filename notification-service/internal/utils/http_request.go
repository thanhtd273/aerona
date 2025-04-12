package utils

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/sony/gobreaker"
)

type HTTPRequestBuilder struct {
	url        string
	method     string
	headers    map[string]string
	body       []byte
	timeout    time.Duration
	maxRetries int
	idempotent bool
	cb         *gobreaker.CircuitBreaker
}

func NewHttpRequestBuilder() *HTTPRequestBuilder {
	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "HTTP client",
		MaxRequests: 2,
		Interval:    60 * time.Second,
		Timeout:     30 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures > 3
		},
	})
	return &HTTPRequestBuilder{
		method:  "GET",
		headers: make(map[string]string),
		timeout: 10 * time.Second,
		cb:      cb,
	}
}

func (b *HTTPRequestBuilder) WithURL(url string) *HTTPRequestBuilder {
	b.url = url
	return b
}

func (b *HTTPRequestBuilder) WithMethod(method string) *HTTPRequestBuilder {
	b.method = method
	return b
}

func (b *HTTPRequestBuilder) WithHeader(key, value string) *HTTPRequestBuilder {
	b.headers[key] = value
	return b
}

func (b *HTTPRequestBuilder) WithBody(body []byte) *HTTPRequestBuilder {
	b.body = body
	return b
}

func (b *HTTPRequestBuilder) WithTimeout(timeout time.Duration) *HTTPRequestBuilder {
	b.timeout = timeout
	return b
}

func (b *HTTPRequestBuilder) WithMaxRetries(maxRetries int) *HTTPRequestBuilder {
	b.maxRetries = maxRetries
	return b
}

func (b *HTTPRequestBuilder) Build() (*http.Response, error) {
	if b.url == "" {
		return nil, fmt.Errorf("url is required")
	}
	if b.method == "" {
		return nil, fmt.Errorf("method is required")
	}

	client := &http.Client{
		Timeout: b.timeout,
	}

	resp, err := b.cb.Execute(func() (any, error) {
		for attempt := 0; attempt <= b.maxRetries; attempt++ {
			req, err := http.NewRequest(b.method, b.url, bytes.NewBuffer(b.body))
			if err != nil {
				return nil, fmt.Errorf("failed to create request: %v", err)
			}
			for k, v := range b.headers {
				req.Header.Add(k, v)
			}
			resp, err := client.Do(req)
			if err != nil {
				if attempt < b.maxRetries && (b.idempotent || attempt == 0) {
					backoff := time.Duration(1<<uint(attempt)) * time.Second
					jitter := time.Duration(rand.Intn(1000)) * time.Millisecond
					time.Sleep(backoff + jitter)
					continue
				}
				return nil, fmt.Errorf("request failed after %d retries: %v", attempt+1, err)
			}

			switch {
			case resp.StatusCode >= 200 && resp.StatusCode < 300:
				return resp, nil
			case resp.StatusCode >= 500 && b.idempotent && attempt < b.maxRetries:
				resp.Body.Close()
				backoff := time.Duration(1<<uint(attempt)) * time.Second
				jitter := time.Duration(rand.Intn(1000)) * time.Millisecond
				time.Sleep(backoff + jitter)
				continue
			default:
				defer resp.Body.Close()
				body, _ := io.ReadAll(resp.Body)
				return nil, fmt.Errorf("request failed with status %d: %v", resp.StatusCode, string(body))
			}
		}
		return nil, errors.New("unexpected retry loop exit")
	})

	if err != nil {
		return nil, fmt.Errorf("circuit breaker or request error: %v", err)
	}

	httpResp, ok := resp.(*http.Response)
	if !ok {
		return nil, errors.New("invalid response type from circuit breaker")
	}
	return httpResp, nil
}
