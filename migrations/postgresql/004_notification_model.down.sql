-- Drop indexes
DROP INDEX IF EXISTS idx_notification_rules_fingerprint;
DROP INDEX IF EXISTS idx_notification_rules_event_level;
DROP INDEX IF EXISTS idx_notification_rules_notification_setting_id;
DROP INDEX IF EXISTS idx_notification_settings_type;
DROP INDEX IF EXISTS idx_notification_settings_project_id;

-- Drop notification_rules table
DROP TABLE IF EXISTS notification_rules;

-- Drop notification_settings table
DROP TABLE IF EXISTS notification_settings;