CREATE TABLE kafka_events (
    event_id String,
    project_id UInt32,
    message String,
    level LowCardinality(String),
    platform LowCardinality(String),
    timestamp DateTime,
    group_hash String,
    source LowCardinality(String),
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
    request_ip Nullable(String),
    -- User
    user_id Nullable(String),
    user_email Nullable(String),
    user_agent Nullable(String),
    -- Contexts
    runtime_name Nullable(String),
    runtime_version Nullable(String),
    os_name Nullable(String),
    os_version Nullable(String),
    browser_name Nullable(String),
    browser_version Nullable(String),
    device_arch Nullable(String),

    server_name String,
    environment String,
    -- Tags
    tags Map(LowCardinality(String), String),
    -- Raw
    raw_data String,
    -- App
    release LowCardinality(Nullable(String))
) ENGINE = Kafka()
SETTINGS kafka_broker_list = 'warden-kafka:9092',
         kafka_topic_list = 'clickhouse.events',
         kafka_group_name = 'clickhouse_events_group',
         kafka_format = 'JSONEachRow';
