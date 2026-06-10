package store

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ervinas/server_uptime_checker/internal/models"
)

// Store handles database reads and writes for monitoring data.
type Store struct {
	db *sql.DB
}

// New creates a Store backed by the given database connection.
func New(db *sql.DB) *Store {
	return &Store{db: db}
}

// ListWebsites returns all websites with their most recent check.
func (s *Store) ListWebsites() ([]models.WebsiteWithLatestCheck, error) {
	rows, err := s.db.Query(`
		SELECT w.id, w.name, w.url, w.created_at, w.last_manual_check_at,
		       c.id, c.status, c.response_time_ms, c.error_message, c.checked_at
		FROM websites w
		LEFT JOIN checks c ON c.id = (
			SELECT id FROM checks
			WHERE website_id = w.id
			ORDER BY checked_at DESC
			LIMIT 1
		)
		ORDER BY w.name ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("list websites: %w", err)
	}
	defer rows.Close()

	var websites []models.WebsiteWithLatestCheck
	for rows.Next() {
		item, err := scanWebsiteWithCheck(rows)
		if err != nil {
			return nil, err
		}
		websites = append(websites, item)
	}
	return websites, rows.Err()
}

// GetWebsite returns one website with its latest check.
func (s *Store) GetWebsite(id int64) (models.WebsiteWithLatestCheck, error) {
	row := s.db.QueryRow(`
		SELECT w.id, w.name, w.url, w.created_at, w.last_manual_check_at,
		       c.id, c.status, c.response_time_ms, c.error_message, c.checked_at
		FROM websites w
		LEFT JOIN checks c ON c.id = (
			SELECT id FROM checks
			WHERE website_id = w.id
			ORDER BY checked_at DESC
			LIMIT 1
		)
		WHERE w.id = ?
	`, id)

	item, err := scanWebsiteWithCheck(row)
	if err == sql.ErrNoRows {
		return models.WebsiteWithLatestCheck{}, fmt.Errorf("website not found")
	}
	if err != nil {
		return models.WebsiteWithLatestCheck{}, fmt.Errorf("get website: %w", err)
	}
	return item, nil
}

// ListChecks returns checks for a website within the last N hours.
func (s *Store) ListChecks(websiteID int64, hours int) ([]models.Check, error) {
	since := time.Now().Add(-time.Duration(hours) * time.Hour)

	rows, err := s.db.Query(`
		SELECT id, website_id, status, response_time_ms, error_message, checked_at
		FROM checks
		WHERE website_id = ? AND checked_at >= ?
		ORDER BY checked_at DESC
	`, websiteID, since)
	if err != nil {
		return nil, fmt.Errorf("list checks: %w", err)
	}
	defer rows.Close()

	var checks []models.Check
	for rows.Next() {
		check, err := scanCheck(rows)
		if err != nil {
			return nil, err
		}
		checks = append(checks, check)
	}
	return checks, rows.Err()
}

// ListAllWebsites returns basic website rows for the worker.
func (s *Store) ListAllWebsites() ([]models.Website, error) {
	rows, err := s.db.Query(`SELECT id, name, url, created_at, last_manual_check_at FROM websites ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("list all websites: %w", err)
	}
	defer rows.Close()

	var websites []models.Website
	for rows.Next() {
		var w models.Website
		var lastManualCheck sql.NullTime
		if err := rows.Scan(&w.ID, &w.Name, &w.URL, &w.CreatedAt, &lastManualCheck); err != nil {
			return nil, err
		}
		if lastManualCheck.Valid {
			w.LastManualCheckAt = &lastManualCheck.Time
		}
		websites = append(websites, w)
	}
	return websites, rows.Err()
}

// InsertCheck saves a probe result and updates incidents when status changes.
func (s *Store) InsertCheck(websiteID int64, result models.CheckStatus, responseTimeMs *int, errorMessage *string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`
		INSERT INTO checks (website_id, status, response_time_ms, error_message, checked_at)
		VALUES (?, ?, ?, ?, NOW())
	`, websiteID, result, responseTimeMs, errorMessage); err != nil {
		return fmt.Errorf("insert check: %w", err)
	}

	var previousStatus sql.NullString
	err = tx.QueryRow(`
		SELECT status FROM checks
		WHERE website_id = ?
		ORDER BY checked_at DESC
		LIMIT 1 OFFSET 1
	`, websiteID).Scan(&previousStatus)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("get previous status: %w", err)
	}

	if err == sql.ErrNoRows || !previousStatus.Valid {
		if result == models.StatusDown {
			if err := openIncident(tx, websiteID); err != nil {
				return err
			}
		}
	} else {
		prev := models.CheckStatus(previousStatus.String)
		if prev == models.StatusUp && result == models.StatusDown {
			if err := openIncident(tx, websiteID); err != nil {
				return err
			}
		}
		if prev == models.StatusDown && result == models.StatusUp {
			if err := closeIncident(tx, websiteID); err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

// DeleteOldChecks removes checks older than the retention period.
func (s *Store) DeleteOldChecks(retentionDays int) error {
	cutoff := time.Now().Add(-time.Duration(retentionDays) * 24 * time.Hour)
	_, err := s.db.Exec(`DELETE FROM checks WHERE checked_at < ?`, cutoff)
	return err
}

// UpdateLastManualCheck updates the timestamp of the last manual probe.
func (s *Store) UpdateLastManualCheck(id int64) error {
	_, err := s.db.Exec(`UPDATE websites SET last_manual_check_at = NOW() WHERE id = ?`, id)
	return err
}

func openIncident(tx *sql.Tx, websiteID int64) error {
	_, err := tx.Exec(`
		INSERT INTO incidents (website_id, started_at)
		VALUES (?, NOW())
	`, websiteID)
	return err
}

func closeIncident(tx *sql.Tx, websiteID int64) error {
	_, err := tx.Exec(`
		UPDATE incidents
		SET ended_at = NOW()
		WHERE website_id = ? AND ended_at IS NULL
	`, websiteID)
	return err
}

type scannable interface {
	Scan(dest ...any) error
}

func scanWebsiteWithCheck(row scannable) (models.WebsiteWithLatestCheck, error) {
	var item models.WebsiteWithLatestCheck
	var checkID sql.NullInt64
	var status sql.NullString
	var responseTime sql.NullInt64
	var errorMessage sql.NullString
	var checkedAt sql.NullTime
	var lastManualCheck sql.NullTime

	err := row.Scan(
		&item.ID, &item.Name, &item.URL, &item.CreatedAt, &lastManualCheck,
		&checkID, &status, &responseTime, &errorMessage, &checkedAt,
	)
	if err != nil {
		return item, err
	}

	if lastManualCheck.Valid {
		item.LastManualCheckAt = &lastManualCheck.Time
	}

	if checkID.Valid {
		check := &models.Check{
			ID:        checkID.Int64,
			WebsiteID: item.ID,
			Status:    models.CheckStatus(status.String),
			CheckedAt: checkedAt.Time,
		}
		if responseTime.Valid {
			rt := int(responseTime.Int64)
			check.ResponseTimeMs = &rt
		}
		if errorMessage.Valid {
			msg := errorMessage.String
			check.ErrorMessage = &msg
		}
		item.LatestCheck = check
	}

	return item, nil
}

func scanCheck(row scannable) (models.Check, error) {
	var check models.Check
	var responseTime sql.NullInt64
	var errorMessage sql.NullString

	err := row.Scan(
		&check.ID, &check.WebsiteID, &check.Status,
		&responseTime, &errorMessage, &check.CheckedAt,
	)
	if err != nil {
		return check, err
	}
	if responseTime.Valid {
		rt := int(responseTime.Int64)
		check.ResponseTimeMs = &rt
	}
	if errorMessage.Valid {
		msg := errorMessage.String
		check.ErrorMessage = &msg
	}
	return check, nil
}
