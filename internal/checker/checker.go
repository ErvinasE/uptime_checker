package checker

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ervinas/server_uptime_checker/internal/models"
)

const defaultTimeout = 10 * time.Second

// Result holds the outcome of a single HTTP probe.
type Result struct {
	Status         models.CheckStatus
	ResponseTimeMs *int
	ErrorMessage   *string
}

// Checker performs HTTP uptime probes.
type Checker struct {
	client *http.Client
}

// New creates a Checker with a 10-second request timeout.
func New() *Checker {
	return &Checker{
		client: &http.Client{
			Timeout: defaultTimeout,
			// Do not follow redirects beyond a reasonable limit.
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= 5 {
					return fmt.Errorf("too many redirects")
				}
				return nil
			},
		},
	}
}

// Check performs a GET request and treats 2xx/3xx as up.
func (c *Checker) Check(ctx context.Context, url string) Result {
	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return downResult(err.Error())
	}
	req.Header.Set("User-Agent", "UptimeMonitor/1.0")

	resp, err := c.client.Do(req)
	elapsed := int(time.Since(start).Milliseconds())

	if err != nil {
		return downResult(err.Error())
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return Result{
			Status:         models.StatusUp,
			ResponseTimeMs: &elapsed,
		}
	}

	msg := fmt.Sprintf("HTTP %d", resp.StatusCode)
	return Result{
		Status:         models.StatusDown,
		ResponseTimeMs: &elapsed,
		ErrorMessage:   &msg,
	}
}

func downResult(message string) Result {
	msg := message
	return Result{
		Status:       models.StatusDown,
		ErrorMessage: &msg,
	}
}
