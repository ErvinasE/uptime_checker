package worker

import (
	"context"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/ervinas/server_uptime_checker/internal/checker"
	"github.com/ervinas/server_uptime_checker/internal/store"
)

const maxConcurrency = 10
const jitterMinutes = 5

// Worker runs periodic uptime checks for all websites.
type Worker struct {
	store            *store.Store
	checker          *checker.Checker
	intervalMinutes  int
	retentionDays    int
}

// New creates a Worker with the given check interval and retention settings.
func New(s *store.Store, c *checker.Checker, intervalMinutes, retentionDays int) *Worker {
	return &Worker{
		store:           s,
		checker:         c,
		intervalMinutes: intervalMinutes,
		retentionDays:   retentionDays,
	}
}

// Run starts the monitoring loop. It runs one check immediately, then repeats
// on a semi-random schedule (base interval ± jitter).
func (w *Worker) Run(ctx context.Context) {
	log.Println("worker: running initial check")
	w.runCycle(ctx)

	for {
		wait := w.nextWait()
		log.Printf("worker: next check in %s", wait)

		select {
		case <-ctx.Done():
			log.Println("worker: stopped")
			return
		case <-time.After(wait):
			w.runCycle(ctx)
		}
	}
}

func (w *Worker) nextWait() time.Duration {
	base := time.Duration(w.intervalMinutes) * time.Minute
	jitter := time.Duration(rand.Intn(jitterMinutes*2+1)-jitterMinutes) * time.Minute
	wait := base + jitter
	if wait < time.Minute {
		wait = time.Minute
	}
	return wait
}

func (w *Worker) runCycle(ctx context.Context) {
	websites, err := w.store.ListAllWebsites()
	if err != nil {
		log.Printf("worker: list websites: %v", err)
		return
	}

	sem := make(chan struct{}, maxConcurrency)
	var wg sync.WaitGroup

	for _, website := range websites {
		wg.Add(1)
		go func(siteID int64, url string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			result := w.checker.Check(ctx, url)
			if err := w.store.InsertCheck(siteID, result.Status, result.ResponseTimeMs, result.ErrorMessage); err != nil {
				log.Printf("worker: save check for site %d: %v", siteID, err)
			}
		}(website.ID, website.URL)
	}

	wg.Wait()

	if err := w.store.DeleteOldChecks(w.retentionDays); err != nil {
		log.Printf("worker: cleanup old checks: %v", err)
	}

	log.Printf("worker: completed check cycle for %d websites", len(websites))
}
