package provider

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"golang.org/x/time/rate"
)

func TestParseRetryAfterSeconds(t *testing.T) {
	got := parseRetryAfter("2", 0)
	if got != 2*time.Second {
		t.Fatalf("parseRetryAfter() = %v, want 2s", got)
	}
}

func TestParseRetryAfterHTTPDate(t *testing.T) {
	retryTime := time.Now().Add(3 * time.Second).UTC()
	got := parseRetryAfter(retryTime.Format(http.TimeFormat), 0)
	if got < 2*time.Second || got > 4*time.Second {
		t.Fatalf("parseRetryAfter() = %v, want about 3s", got)
	}
}

func TestParseRetryAfterMissingUsesBackoff(t *testing.T) {
	tests := []struct {
		attempt int
		want    time.Duration
	}{
		{0, 1 * time.Second},
		{1, 2 * time.Second},
		{2, 4 * time.Second},
	}

	for _, tt := range tests {
		got := parseRetryAfter("", tt.attempt)
		if got != tt.want {
			t.Fatalf("parseRetryAfter('', %d) = %v, want %v", tt.attempt, got, tt.want)
		}
	}
}

func TestThrottledTransportRetriesAfter429WithRetryAfterHeader(t *testing.T) {
	var requestCount atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if requestCount.Add(1) == 1 {
			w.Header().Set("Retry-After", "1")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	transport := &ThrottledTransport{
		Transport: http.DefaultTransport,
		limiter:   rate.NewLimiter(rate.Inf, 1),
	}
	client := &http.Client{Transport: transport}

	start := time.Now()
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL, nil)
	if err != nil {
		t.Fatalf("NewRequestWithContext() error = %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
	if requestCount.Load() != 2 {
		t.Fatalf("request count = %d, want 2", requestCount.Load())
	}
	if elapsed := time.Since(start); elapsed < time.Second {
		t.Fatalf("elapsed = %v, want at least 1s", elapsed)
	}
}

func TestThrottledTransportGlobalPauseOn429(t *testing.T) {
	var (
		requestCount atomic.Int32
		requestMu    sync.Mutex
		requestTimes []time.Time
	)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestMu.Lock()
		requestTimes = append(requestTimes, time.Now())
		requestMu.Unlock()

		if requestCount.Add(1) <= 2 {
			w.Header().Set("Retry-After", "1")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	transport := &ThrottledTransport{
		Transport: http.DefaultTransport,
		limiter:   rate.NewLimiter(rate.Inf, 1),
	}
	client := &http.Client{Transport: transport}

	var wg sync.WaitGroup
	wg.Add(2)

	for i := 0; i < 2; i++ {
		go func() {
			defer wg.Done()

			req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL, nil)
			if err != nil {
				t.Errorf("NewRequestWithContext() error = %v", err)
				return
			}

			resp, err := client.Do(req)
			if err != nil {
				t.Errorf("Do() error = %v", err)
				return
			}
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusOK)
			}
		}()
	}

	wg.Wait()

	if requestCount.Load() != 4 {
		t.Fatalf("request count = %d, want 4", requestCount.Load())
	}

	requestMu.Lock()
	defer requestMu.Unlock()

	if len(requestTimes) < 4 {
		t.Fatalf("got %d request timestamps, want at least 4", len(requestTimes))
	}

	pauseWindowStart := requestTimes[1]
	pauseWindowEnd := requestTimes[2]
	if pauseWindowEnd.Sub(pauseWindowStart) < 900*time.Millisecond {
		t.Fatalf("requests during pause were too close together: %v", pauseWindowEnd.Sub(pauseWindowStart))
	}
}

func TestThrottledTransportPersistent429ReturnsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "0")
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	transport := &ThrottledTransport{
		Transport: http.DefaultTransport,
		limiter:   rate.NewLimiter(rate.Inf, 1),
	}
	client := &http.Client{Transport: transport}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL, nil)
	if err != nil {
		t.Fatalf("NewRequestWithContext() error = %v", err)
	}

	_, err = client.Do(req)
	if err == nil {
		t.Fatal("Do() error = nil, want rate limit error")
	}
}

func TestThrottledTransportLimiterSmoke(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	transport := &ThrottledTransport{
		Transport: http.DefaultTransport,
		limiter:   rate.NewLimiter(rate.Limit(5), 1),
	}
	client := &http.Client{Transport: transport}

	start := time.Now()
	for i := 0; i < 5; i++ {
		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL, nil)
		if err != nil {
			t.Fatalf("NewRequestWithContext() error = %v", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Do() error = %v", err)
		}
		resp.Body.Close()
	}

	elapsed := time.Since(start)
	if elapsed < 600*time.Millisecond {
		t.Fatalf("elapsed = %v, expected rate limiting to slow requests", elapsed)
	}
}

func TestNewAllQuietAPIClientUsesThrottledTransport(t *testing.T) {
	client := NewAllQuietAPIClient("test-key", "https://example.com", nil, "test")
	transport, ok := client.HTTPClient.Transport.(*ThrottledTransport)
	if !ok {
		t.Fatalf("transport type = %T, want *ThrottledTransport", client.HTTPClient.Transport)
	}

	authTransport, ok := transport.Transport.(*AuthTransport)
	if !ok {
		t.Fatalf("inner transport type = %T, want *AuthTransport", transport.Transport)
	}

	if authTransport.APIKey != "test-key" {
		t.Fatalf("API key = %q, want test-key", authTransport.APIKey)
	}
}
