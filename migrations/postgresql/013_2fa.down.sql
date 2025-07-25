ALTER TABLE users
    DROP COLUMN two_fa_enabled,
    DROP COLUMN two_fa_secret,
    DROP COLUMN two_fa_confirmed_at;
