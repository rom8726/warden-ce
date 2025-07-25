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