CREATE TABLE notifications_queue (
    id BIGSERIAL PRIMARY KEY,
    issue_id INTEGER NOT NULL REFERENCES issues(id) ON DELETE CASCADE,
    project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    level TEXT NOT NULL,
    is_new BOOLEAN DEFAULT true,
    was_reactivated BOOLEAN DEFAULT false,
    sent_at TIMESTAMPTZ,
    status TEXT CHECK (status IN ('pending', 'sent', 'failed', 'skipped')) DEFAULT 'pending',
    fail_reason TEXT,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_notifications_queue_status ON notifications_queue(status);
CREATE INDEX idx_notifications_queue_project_issue ON notifications_queue(project_id, issue_id);
