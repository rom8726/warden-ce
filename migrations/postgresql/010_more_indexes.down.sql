-- Rollback for additional indexes
-- Drop in reverse order to the creation to avoid dependencies

-- ====== USERS ==============================================================
DROP INDEX IF EXISTS idx_users_active_created_at;

-- ====== TEAM_MEMBERS =======================================================
DROP INDEX IF EXISTS idx_tm_user_role;

-- ====== NOTIFICATION_SETTINGS =============================================
DROP INDEX IF EXISTS idx_ns_config_gin;
DROP INDEX IF EXISTS idx_ns_project_enabled;

-- ====== NOTIFICATIONS_QUEUE ===============================================
DROP INDEX IF EXISTS idx_nq_issue_id;
DROP INDEX IF EXISTS idx_nq_project_sent_at;
DROP INDEX IF EXISTS idx_nq_pending;

-- ====== RESOLUTIONS ========================================================
DROP INDEX IF EXISTS idx_resolutions_open;
DROP INDEX IF EXISTS idx_resolutions_project_status;

-- ====== ISSUES =============================================================
DROP INDEX IF EXISTS idx_issues_title_tsv;
DROP INDEX IF EXISTS idx_issues_fingerprint;
DROP INDEX IF EXISTS idx_issues_unresolved;
DROP INDEX IF EXISTS idx_issues_project_status_last_seen;
