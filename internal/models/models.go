package models

import "time"

// Website is a monitored URL.
type Website struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
	LastManualCheckAt *time.Time `json:"last_manual_check_at,omitempty"`
}

// CheckStatus is the result of a single probe.
type CheckStatus string

const (
	StatusUp   CheckStatus = "up"
	StatusDown CheckStatus = "down"
)

// Check records one uptime probe result.
type Check struct {
	ID             int64       `json:"id"`
	WebsiteID      int64       `json:"website_id"`
	Status         CheckStatus `json:"status"`
	ResponseTimeMs *int        `json:"response_time_ms"`
	ErrorMessage   *string     `json:"error_message,omitempty"`
	CheckedAt      time.Time   `json:"checked_at"`
}

// Incident tracks a downtime period for a website.
type Incident struct {
	ID        int64      `json:"id"`
	WebsiteID int64      `json:"website_id"`
	StartedAt time.Time  `json:"started_at"`
	EndedAt   *time.Time `json:"ended_at,omitempty"`
}

// WebsiteWithLatestCheck is returned by list/detail API endpoints.
type WebsiteWithLatestCheck struct {
	Website
	LatestCheck *Check `json:"latest_check,omitempty"`
}
