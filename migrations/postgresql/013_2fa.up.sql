ALTER TABLE users
    ADD COLUMN two_fa_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN two_fa_secret TEXT,
    ADD COLUMN two_fa_confirmed_at TIMESTAMPTZ;
