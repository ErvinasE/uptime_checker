-- Add last_manual_check_at to websites table
ALTER TABLE websites ADD COLUMN last_manual_check_at TIMESTAMP;
