-- Materialized view to transfer data from Kafka to the events table
CREATE MATERIALIZED VIEW mv_events TO events AS
SELECT
    event_id,
    project_id,
    message,
    level,
    platform,
    timestamp,
    group_hash,
    source,

    stacktrace,
    exception_type,
    exception_value,

    request_url,
    request_method,
    request_query,
    request_headers,
    request_data,
    request_cookies,
    request_ip,

    user_id,
    user_email,
    user_agent,

    runtime_name,
    runtime_version,
    os_name,
    os_version,
    browser_name,
    browser_version,
    device_arch,

    server_name,
    environment,
    tags,
    raw_data,

    release
FROM kafka_events;
