CREATE TYPE issue_source AS ENUM ('event', 'exception');
CREATE TYPE issue_status AS ENUM ('unresolved', 'resolved', 'ignored');

CREATE TABLE issues (
    id BIGSERIAL PRIMARY KEY,
    project_id INTEGER NOT NULL REFERENCES projects(id),
    fingerprint TEXT NOT NULL,
    source issue_source NOT NULL DEFAULT 'event',
    status issue_status NOT NULL DEFAULT 'unresolved',
    title TEXT,
    level TEXT,
    platform TEXT,
    first_seen TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_seen TIMESTAMPTZ NOT NULL DEFAULT now(),
    total_events INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    UNIQUE (project_id, fingerprint)
);

CREATE INDEX idx_issues_project_last_seen ON issues(project_id, last_seen DESC);
