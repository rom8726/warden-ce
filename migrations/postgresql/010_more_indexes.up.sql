-- ====== ISSUES =============================================================

-- 1. Most frequently used: list “live” issues for a project sorted by time.
CREATE INDEX IF NOT EXISTS idx_issues_project_status_last_seen
    ON issues (project_id, status, last_seen DESC);

-- 2. Fast lookup of “unresolved” issues.
CREATE INDEX IF NOT EXISTS idx_issues_unresolved
    ON issues (project_id, last_seen DESC)
    WHERE status = 'unresolved';

-- 3. Search by fingerprint outside of project context.
CREATE INDEX IF NOT EXISTS idx_issues_fingerprint
    ON issues (fingerprint);

-- 4. Full-text search by title.
CREATE INDEX IF NOT EXISTS idx_issues_title_tsv
    ON issues USING GIN (to_tsvector('simple', title));

-- ====== RESOLUTIONS ========================================================

-- 5. Project/status statistics for resolutions.
CREATE INDEX IF NOT EXISTS idx_resolutions_project_status
    ON resolutions (project_id, status);

-- 6. Open resolutions only (partial index).
CREATE INDEX IF NOT EXISTS idx_resolutions_open
    ON resolutions (project_id)
    WHERE status <> 'resolved';

-- ====== NOTIFICATIONS_QUEUE ===============================================

-- 7. Worker: scan “pending” items to send.
CREATE INDEX IF NOT EXISTS idx_nq_pending
    ON notifications_queue (project_id, issue_id)
    WHERE status = 'pending';

-- 8. History of sent notifications sorted by time.
CREATE INDEX IF NOT EXISTS idx_nq_project_sent_at
    ON notifications_queue (project_id, sent_at DESC)
    WHERE status = 'sent';

-- 9. Quick jump from issue to queue items.
CREATE INDEX IF NOT EXISTS idx_nq_issue_id
    ON notifications_queue (issue_id);

-- ====== NOTIFICATION_SETTINGS =============================================

-- 10. Active settings of a project.
CREATE INDEX IF NOT EXISTS idx_ns_project_enabled
    ON notification_settings (project_id)
    WHERE enabled = TRUE;

-- 11. GIN index for JSON config lookups.
CREATE INDEX IF NOT EXISTS idx_ns_config_gin
    ON notification_settings USING GIN (config);

-- ====== TEAM_MEMBERS =======================================================

-- 12. Get user role inside teams.
CREATE INDEX IF NOT EXISTS idx_tm_user_role
    ON team_members (user_id, role);
