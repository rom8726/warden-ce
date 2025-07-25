-- Create resolutions table
CREATE TABLE IF NOT EXISTS resolutions (
    id SERIAL PRIMARY KEY,
    project_id INT NOT NULL REFERENCES projects(id),
    issue_id INT NOT NULL REFERENCES issues(id),
    status issue_status NOT NULL DEFAULT 'unresolved',
    resolved_by INT REFERENCES users(id),
    resolved_at TIMESTAMP,
    comment TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for faster lookups
CREATE INDEX IF NOT EXISTS idx_resolutions_project_id ON resolutions(project_id);
CREATE INDEX IF NOT EXISTS idx_resolutions_issue_id ON resolutions(issue_id);
CREATE INDEX IF NOT EXISTS idx_resolutions_status ON resolutions(status);
CREATE INDEX IF NOT EXISTS idx_resolutions_resolved_by ON resolutions(resolved_by);