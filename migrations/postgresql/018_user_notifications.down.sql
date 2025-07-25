-- Drop indexes
DROP INDEX IF EXISTS idx_user_notifications_user_read;
DROP INDEX IF EXISTS idx_user_notifications_created_at;
DROP INDEX IF EXISTS idx_user_notifications_is_read;
DROP INDEX IF EXISTS idx_user_notifications_user_id;

-- Drop table
DROP TABLE IF EXISTS user_notifications;
