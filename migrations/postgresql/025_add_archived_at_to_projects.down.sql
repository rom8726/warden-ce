-- Drop the index first
DROP INDEX IF EXISTS idx_projects_archived_at;

-- Remove the archived_at column
ALTER TABLE projects DROP COLUMN IF EXISTS archived_at;