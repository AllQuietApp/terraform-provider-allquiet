package provider

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golang.org/x/time/rate"
)

const (
	defaultRequestsPerSecond = rate.Limit(8)
	defaultBurst             = 10
	maxRateLimitRetries      = 5
	maxRetryAfterSeconds     = 300
)

type ThrottledTransport struct {
	Transport http.RoundTripper
	limiter   *rate.Limiter

	mu         sync.Mutex
	pauseUntil time.Time
}

func newThrottledTransport(inner http.RoundTripper) *ThrottledTransport {
	if inner == nil {
		inner = http.DefaultTransport
	}

	return &ThrottledTransport{
		Transport: inner,
		limiter:   rate.NewLimiter(defaultRequestsPerSecond, defaultBurst),
	}
}

func (t *ThrottledTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var lastErr error

	for attempt := 0; attempt < maxRateLimitRetries; attempt++ {
		if err := t.limiter.Wait(req.Context()); err != nil {
			return nil, err
		}

		if err := t.waitForPause(req.Context()); err != nil {
			return nil, err
		}

		resp, err := t.Transport.RoundTrip(req)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != http.StatusTooManyRequests {
			return resp, nil
		}

		retryAfter := parseRetryAfter(resp.Header.Get("Retry-After"), attempt)
		resp.Body.Close()

		t.extendPause(retryAfter)
		tflog.Warn(req.Context(), "Rate limited by All Quiet API, retrying after pause", map[string]interface{}{
			"retry_after_seconds": retryAfter.Seconds(),
			"attempt":             attempt + 1,
		})

		if attempt == maxRateLimitRetries-1 {
			lastErr = fmt.Errorf("All Quiet API rate limit exceeded after %d retries", maxRateLimitRetries)
			break
		}

		if err := t.waitForPause(req.Context()); err != nil {
			return nil, err
		}

		if req.GetBody != nil {
			body, err := req.GetBody()
			if err != nil {
				return nil, err
			}
			req.Body = body
		}
	}

	return nil, lastErr
}

func (t *ThrottledTransport) extendPause(duration time.Duration) {
	t.mu.Lock()
	defer t.mu.Unlock()

	until := time.Now().Add(duration)
	if until.After(t.pauseUntil) {
		t.pauseUntil = until
	}
}

func (t *ThrottledTransport) waitForPause(ctx context.Context) error {
	t.mu.Lock()
	until := t.pauseUntil
	t.mu.Unlock()

	wait := time.Until(until)
	if wait <= 0 {
		return nil
	}

	tflog.Debug(ctx, "Waiting for All Quiet API rate limit pause", map[string]interface{}{
		"wait_seconds": wait.Seconds(),
	})

	timer := time.NewTimer(wait)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func parseRetryAfter(header string, attempt int) time.Duration {
	if header == "" {
		return fallbackRetryAfter(attempt)
	}

	if seconds, err := strconv.Atoi(header); err == nil {
		if seconds <= 0 {
			return fallbackRetryAfter(attempt)
		}
		if seconds > maxRetryAfterSeconds {
			seconds = maxRetryAfterSeconds
		}
		return time.Duration(seconds) * time.Second
	}

	if retryTime, err := http.ParseTime(header); err == nil {
		wait := time.Until(retryTime)
		if wait <= 0 {
			return fallbackRetryAfter(attempt)
		}
		if wait > maxRetryAfterSeconds*time.Second {
			wait = maxRetryAfterSeconds * time.Second
		}
		return wait
	}

	return fallbackRetryAfter(attempt)
}

func fallbackRetryAfter(attempt int) time.Duration {
	wait := time.Duration(1<<attempt) * time.Second
	maxWait := maxRetryAfterSeconds * time.Second
	if wait > maxWait {
		wait = maxWait
	}
	return wait
}
