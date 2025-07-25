CREATE TABLE events (
    event_id String,
    project_id UInt32,
    message String,
    level LowCardinality(String),   -- error/warning/info/exception
    platform LowCardinality(String),
    timestamp DateTime DEFAULT now(),
    group_hash String,
    source LowCardinality(String),  -- event/exception
    server_name String,
    environment String,
    tags Map(LowCardinality(String), String) DEFAULT map(),
    -- Exception
    stacktrace Nullable(String),
    exception_type Nullable(String),
    exception_value Nullable(String),
    -- Request context
    request_url Nullable(String),
    request_method Nullable(String),
    request_query Nullable(String),
    request_headers Map(LowCardinality(String), String),
    request_data Nullable(String),
    request_cookies Nullable(String),
    request_ip Nullable(String),  -- REMOTE_ADDR or user.ip_address
    -- User
    user_id Nullable(String),
    user_email Nullable(String),
    user_agent Nullable(String), -- from headers["User-Agent"]
    -- Contexts
    runtime_name Nullable(String),
    runtime_version Nullable(String),
    os_name Nullable(String),
    os_version Nullable(String),
    browser_name Nullable(String),
    browser_version Nullable(String),
    device_arch Nullable(String),
    -- Raw data
    raw_data String,
    -- App
    release LowCardinality(Nullable(String))
) ENGINE = MergeTree()
PARTITION BY toYYYYMM(timestamp)
ORDER BY (project_id, timestamp, group_hash)
TTL timestamp + INTERVAL 3 MONTH
SETTINGS index_granularity = 8192;
