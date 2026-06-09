package checker

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ervinas/server_uptime_checker/internal/models"
)

const defaultTimeout = 20 * time.Second

// browserUserAgent mimics a normal browser so bot filters are less likely to block us.
const browserUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"

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

// New creates a Checker with a 20-second request timeout.
func New() *Checker {
	return &Checker{
		client: &http.Client{
			Timeout: defaultTimeout,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= 5 {
					return fmt.Errorf("too many redirects")
				}
				return nil
			},
		},
	}
}

// Check performs a GET request. The site is considered up when the server
// responds with any HTTP status below 500. Many popular sites return 403 to
// automated clients even when they are online for real users.
func (c *Checker) Check(ctx context.Context, url string) Result {
	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return downResult(err.Error())
	}
	req.Header.Set("User-Agent", browserUserAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	resp, err := c.client.Do(req)
	elapsed := int(time.Since(start).Milliseconds())

	if err != nil {
		return downResult(err.Error())
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)

	if resp.StatusCode < 500 {
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
