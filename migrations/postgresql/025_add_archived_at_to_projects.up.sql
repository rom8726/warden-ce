-- Add archived_at column to projects table
ALTER TABLE projects ADD COLUMN archived_at TIMESTAMP NULL;

-- Create index for faster lookups of non-archived projects
CREATE INDEX idx_projects_archived_at ON projects(archived_at) WHERE archived_at IS NULL;