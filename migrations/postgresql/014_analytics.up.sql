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
