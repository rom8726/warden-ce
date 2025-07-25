-- Create projects table
CREATE TABLE IF NOT EXISTS projects (
                                        id SERIAL PRIMARY KEY,
                                        name VARCHAR(255) NOT NULL,
                                        public_key VARCHAR(255) NOT NULL UNIQUE,
                                        created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for faster lookups
CREATE INDEX IF NOT EXISTS idx_projects_public_key ON projects(public_key);

-- Create users table
CREATE TABLE IF NOT EXISTS users (
                                     id SERIAL PRIMARY KEY,
                                     username VARCHAR(255) NOT NULL UNIQUE,
                                     email VARCHAR(255) NOT NULL UNIQUE,
                                     password_hash VARCHAR(255) NOT NULL,
                                     is_superuser BOOLEAN NOT NULL DEFAULT FALSE,
                                     is_active BOOLEAN NOT NULL DEFAULT TRUE,
                                     created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                     updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                     last_login TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for faster lookups
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Create teams table
CREATE TABLE IF NOT EXISTS teams (
                                     id SERIAL PRIMARY KEY,
                                     name VARCHAR(255) NOT NULL,
                                     created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create team_members table
CREATE TABLE IF NOT EXISTS team_members (
                                            team_id INT NOT NULL REFERENCES teams(id),
                                            user_id INT NOT NULL REFERENCES users(id),
                                            role VARCHAR(50) NOT NULL, -- owner, admin, member
                                            PRIMARY KEY (team_id, user_id)
);

-- Create indexes for faster lookups
CREATE INDEX IF NOT EXISTS idx_teams_name ON teams(name);
CREATE INDEX IF NOT EXISTS idx_team_members_user_id ON team_members(user_id);

-- Add team_id column to projects table
ALTER TABLE projects ADD COLUMN team_id INT REFERENCES teams(id);
CREATE INDEX IF NOT EXISTS idx_projects_team_id ON projects(team_id);

-- Create notification_settings table
CREATE TABLE IF NOT EXISTS notification_settings (
                                                     id SERIAL PRIMARY KEY,
                                                     project_id INT NOT NULL REFERENCES projects(id),
                                                     type VARCHAR(50) NOT NULL, -- email, telegram, slack
                                                     config JSONB NOT NULL,
                                                     enabled BOOLEAN NOT NULL DEFAULT TRUE,
                                                     created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                                                     updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create notification_rules table
CREATE TABLE IF NOT EXISTS notification_rules (
                                                  id SERIAL PRIMARY KEY,
                                                  notification_setting_id INT NOT NULL REFERENCES notification_settings(id),
                                                  event_level VARCHAR(50), -- error, warning, info, etc.
                                                  fingerprint TEXT, -- specific error fingerprint
                                                  is_new_error BOOLEAN, -- trigger only for new errors
                                                  is_regression BOOLEAN, -- trigger only for regressions (resolved -> unresolved)
                                                  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for faster lookups
CREATE INDEX IF NOT EXISTS idx_notification_settings_project_id ON notification_settings(project_id);
CREATE INDEX IF NOT EXISTS idx_notification_settings_type ON notification_settings(type);
CREATE INDEX IF NOT EXISTS idx_notification_rules_notification_setting_id ON notification_rules(notification_setting_id);
CREATE INDEX IF NOT EXISTS idx_notification_rules_event_level ON notification_rules(event_level);
CREATE INDEX IF NOT EXISTS idx_notification_rules_fingerprint ON notification_rules(fingerprint);

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

-- issue fingerprint is 40 chars

ALTER TABLE issues ALTER COLUMN fingerprint TYPE CHAR(40);

-- users have tmp passwords

ALTER TABLE users ADD COLUMN is_tmp_password boolean DEFAULT true;

-- Create the notifications_queue table

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

-- issues.last_notification_at

ALTER TABLE issues ADD COLUMN last_notification_at timestamp with time zone;

-- projects.description

ALTER TABLE projects ADD COLUMN description TEXT;

-- users.2fa columns

ALTER TABLE users
    ADD COLUMN two_fa_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN two_fa_secret TEXT,
    ADD COLUMN two_fa_confirmed_at TIMESTAMPTZ;

-- Analytics

CREATE TABLE releases (
                          id SERIAL PRIMARY KEY,
                          project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
                          version TEXT NOT NULL,
                          description TEXT,
                          released_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                          created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                          updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

                          UNIQUE(project_id, version)
);

CREATE TABLE issue_releases (
                                issue_id BIGINT NOT NULL REFERENCES issues(id) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED,
                                release_id INTEGER NOT NULL REFERENCES releases(id) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED,
                                first_seen_in BOOLEAN DEFAULT FALSE, -- issue впервые замечено в этом релизе?
                                PRIMARY KEY (issue_id, release_id)
);

CREATE TABLE release_stats (
                               project_id                   INT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
                               release_id                   INTEGER NOT NULL REFERENCES releases(id) ON DELETE CASCADE,
                               release                      TEXT NOT NULL,
                               generated_at                 TIMESTAMPTZ NOT NULL DEFAULT now(),

                               known_issues_total           INT NOT NULL DEFAULT 0, -- общее кол-во известных на момент генерации
                               new_issues_total             INT NOT NULL DEFAULT 0, -- впервые появившихся в этой версии
                               regressions_total            INT NOT NULL DEFAULT 0, -- регрессий

                               resolved_in_version_total    INT NOT NULL DEFAULT 0, -- все решенные в этой версии
                               fixed_new_in_version_total   INT NOT NULL DEFAULT 0, -- новые, появившиеся и решенные в этой версии
                               fixed_old_in_version_total   INT NOT NULL DEFAULT 0, -- старые, решенные в этой версии

                               avg_fix_time                 TEXT NOT NULL DEFAULT '0s', -- среднее время фикса
                               median_fix_time              TEXT NOT NULL DEFAULT '0s', -- 50-й перцентиль
                               p95_fix_time                 TEXT NOT NULL DEFAULT '0s', -- 95-й перцентиль

                               severity_distribution        JSONB, -- JSONB-колонка: распределение ошибок по level

                               users_affected               INT, -- сколько уникальных user_id было в event'ах этой версии, будет считаться по ClickHouse: count(distinct user_id).

                               PRIMARY KEY (project_id, release)
);

-- Create a settings table for storing application configuration
CREATE TABLE settings (
                          id SERIAL PRIMARY KEY,
                          name VARCHAR(255) NOT NULL UNIQUE,
                          value JSONB NOT NULL,
                          description TEXT,
                          created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                          updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create an index on name for fast lookups
CREATE INDEX idx_settings_name ON settings(name);

-- Create trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_settings_updated_at
    BEFORE UPDATE ON settings
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Create user_notifications table
CREATE TABLE user_notifications (
                                    id SERIAL PRIMARY KEY,
                                    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                    type VARCHAR(50) NOT NULL,
                                    content JSONB NOT NULL,
                                    is_read BOOLEAN NOT NULL DEFAULT FALSE,
                                    email_sent BOOLEAN NOT NULL DEFAULT FALSE,
                                    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create indexes for performance
CREATE INDEX idx_user_notifications_user_id ON user_notifications(user_id);
CREATE INDEX idx_user_notifications_is_read ON user_notifications(is_read);
CREATE INDEX idx_user_notifications_created_at ON user_notifications(created_at);
CREATE INDEX idx_user_notifications_user_read ON user_notifications(user_id, is_read);

-- Add archived_at column to projects table
ALTER TABLE projects ADD COLUMN archived_at TIMESTAMP NULL;

-- Create index for faster lookups of non-archived projects
CREATE INDEX idx_projects_archived_at ON projects(archived_at) WHERE archived_at IS NULL;
