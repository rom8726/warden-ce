package issues

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	wardencontext "github.com/rom8726/warden/internal/context"
	"github.com/rom8726/warden/internal/domain"
	"github.com/rom8726/warden/pkg/db"
)

type Repository struct {
	db db.Tx
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{
		db: pool,
	}
}

func (r *Repository) GetByID(ctx context.Context, id domain.IssueID) (domain.Issue, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM issues WHERE id = $1 LIMIT 1`
	rows, err := executor.Query(ctx, query, id)
	if err != nil {
		return domain.Issue{}, fmt.Errorf("query issue by ID: %w", err)
	}
	defer rows.Close()

	issue, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[issueModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Issue{}, domain.ErrEntityNotFound
		}

		return domain.Issue{}, fmt.Errorf("collect user: %w", err)
	}

	return issue.toDomain(), nil
}

func (r *Repository) ListByFingerprints(ctx context.Context, fingerprints []string) ([]domain.Issue, error) {
	executor := r.getExecutor(ctx)
	const query = `SELECT * FROM issues WHERE fingerprint = ANY($1) LIMIT $2`
	rows, err := executor.Query(ctx, query, fingerprints, len(fingerprints))
	if err != nil {
		return nil, fmt.Errorf("query issues by fingerprints: %w", err)
	}
	defer rows.Close()

	list, err := pgx.CollectRows(rows, pgx.RowToStructByName[issueModel])
	if err != nil {
		return nil, fmt.Errorf("collect issues by fingerprints: %w", err)
	}

	issues := make([]domain.Issue, 0, len(list))
	for _, issue := range list {
		issues = append(issues, issue.toDomain())
	}

	return issues, nil
}

func (r *Repository) UpsertIssue(ctx context.Context, issue domain.IssueDTO) (domain.IssueUpsertResult, error) {
	executor := r.getExecutor(ctx)

	type UpsertResult struct {
		ID             uint  `db:"id"`
		IsNew          bool  `db:"is_new"`
		WasReactivated *bool `db:"was_reactivated"`
	}

	var res UpsertResult

	query := `
WITH old AS (
  SELECT status
  FROM issues
  WHERE project_id = $1 AND fingerprint = $2
),
upserted AS (
  INSERT INTO issues (
    project_id, fingerprint, source, status,
    title, level, platform,
    first_seen, last_seen, total_events
  )
  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $8, 1)
  ON CONFLICT (project_id, fingerprint)
  DO UPDATE SET
    last_seen = GREATEST(issues.last_seen, EXCLUDED.last_seen),
    status = CASE
               WHEN issues.status = 'ignored' THEN 'ignored'
               ELSE EXCLUDED.status
             END,
    total_events = issues.total_events + 1
  RETURNING id, first_seen, last_seen
)
SELECT
  upserted.id,
  (first_seen = last_seen) AS is_new,
  COALESCE(old.status = 'resolved', false) AS was_reactivated
FROM upserted
LEFT JOIN old ON TRUE;`

	err := executor.QueryRow(ctx, query,
		issue.ProjectID,
		issue.Fingerprint,
		issue.Source,
		domain.IssueStatusUnresolved,
		issue.Title,
		issue.Level,
		issue.Platform,
		time.Now(),
	).Scan(&res.ID, &res.IsNew, &res.WasReactivated)

	var wasReactivated bool
	if res.WasReactivated != nil && *res.WasReactivated {
		wasReactivated = true
	}

	return domain.IssueUpsertResult{
		ID:             domain.IssueID(res.ID),
		IsNew:          res.IsNew,
		WasReactivated: wasReactivated,
	}, err
}

//nolint:gocyclo // need refactoring
func (r *Repository) ListExtended(
	ctx context.Context,
	filter *domain.ListIssuesFilter,
) ([]domain.IssueExtended, uint64, error) {
	executor := r.getExecutor(ctx)

	uid := wardencontext.UserID(ctx)
	isSuper := wardencontext.IsSuper(ctx)

	page := filter.PageNum
	if page <= 0 {
		page = 1
	}
	limit := filter.PerPage
	if limit <= 0 {
		limit = 20
	}
	offset := (page - 1) * limit

	stmt := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	baseSelect := stmt.
		Select(
			"issues.id AS id",
			"issues.project_id AS project_id",
			"projects.name AS project_name",
			"fingerprint",
			"source",
			"issues.status AS status",
			"title",
			"level",
			"platform",
			"first_seen",
			"last_seen",
			"total_events",
			"issues.created_at AS created_at",
			"issues.updated_at AS updated_at",
			"r.resolved_by",
			"r.resolved_at",
			"resolver.username AS resolved_by_username",
		).
		From("issues").
		LeftJoin("projects ON issues.project_id = projects.id").
		LeftJoin(`LATERAL (
			SELECT *
			FROM resolutions
			WHERE resolutions.issue_id = issues.id
			ORDER BY resolved_at DESC
			LIMIT 1
		) r ON true`).
		LeftJoin("users resolver ON resolver.id = r.resolved_by").
		LeftJoin("team_members tm ON tm.team_id = projects.team_id AND tm.user_id = ?", uid).
		Where("projects.archived_at IS NULL")

	baseCount := stmt.Select("COUNT(*)").
		From("issues").
		Join("projects ON projects.id = issues.project_id").
		LeftJoin("team_members tm ON tm.team_id = projects.team_id AND tm.user_id = ?", uid).
		Where("projects.archived_at IS NULL")

	// permission condition for non-superusers
	if !isSuper {
		permCond := sq.Or{
			sq.Expr("projects.team_id IS NULL"),
			sq.Expr("tm.user_id IS NOT NULL"),
		}
		baseSelect = baseSelect.Where(permCond)
		baseCount = baseCount.Where(permCond)
	}

	if filter.ProjectID != nil {
		baseSelect = baseSelect.Where("issues.project_id = ?", *filter.ProjectID)
		baseCount = baseCount.Where("issues.project_id = ?", *filter.ProjectID)
	}

	if filter.Level != nil {
		baseSelect = baseSelect.Where("issues.level = ?", *filter.Level)
		baseCount = baseCount.Where("issues.level = ?", *filter.Level)
	}

	if filter.Status != nil {
		baseSelect = baseSelect.Where("issues.status = ?", *filter.Status)
		baseCount = baseCount.Where("issues.status = ?", *filter.Status)
	}

	if !filter.TimeFrom.IsZero() {
		baseSelect = baseSelect.Where("issues.last_seen >= ?", filter.TimeFrom)
		baseCount = baseCount.Where("issues.last_seen >= ?", filter.TimeFrom)
	}

	if !filter.TimeTo.IsZero() {
		baseSelect = baseSelect.Where("issues.last_seen <= ?", filter.TimeTo)
		baseCount = baseCount.Where("issues.last_seen <= ?", filter.TimeTo)
	}

	// Build the order by clause based on the filter
	orderByBuilder := baseSelect.OrderBy(
		`CASE level 
			WHEN 'fatal' THEN 1 
			WHEN 'exception' THEN 2 
			WHEN 'error' THEN 3 
			WHEN 'warning' THEN 4 
			WHEN 'info' THEN 5 
			WHEN 'debug' THEN 6 
			ELSE 7 
		END`,
	)

	// Add the sort column based on the filter
	var sortColumn string
	switch filter.OrderBy {
	case domain.OrderByFieldTotalEvents:
		sortColumn = "issues.total_events"
	case domain.OrderByFieldFirstSeen:
		sortColumn = "issues.first_seen"
	case domain.OrderByFieldLastSeen:
		sortColumn = "issues.last_seen"
	default:
		// Default to total_events if not specified
		sortColumn = "issues.total_events"
	}

	// Add the sort direction
	if filter.OrderAsc {
		sortColumn += " ASC"
	} else {
		sortColumn += " DESC"
	}

	orderByBuilder = orderByBuilder.OrderBy(sortColumn)

	selectSQL, selectArgs, err := orderByBuilder.
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		ToSql()
	if err != nil {
		return nil, 0, err
	}

	countSQL, countArgs, err := baseCount.ToSql()
	if err != nil {
		return nil, 0, err
	}

	var total uint64
	if err := executor.QueryRow(ctx, countSQL, countArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := executor.Query(ctx, selectSQL, selectArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	issues := make([]domain.IssueExtended, 0, limit)

	for rows.Next() {
		var is domain.IssueExtended
		var resolvedBy pgtype.Int4
		var resolvedAt pgtype.Timestamptz
		var resolvedByUsername pgtype.Text

		if err := rows.Scan(
			&is.ID,
			&is.ProjectID,
			&is.ProjectName,
			&is.Fingerprint,
			&is.Source,
			&is.Status,
			&is.Title,
			&is.Level,
			&is.Platform,
			&is.FirstSeen,
			&is.LastSeen,
			&is.TotalEvents,
			&is.CreatedAt,
			&is.UpdatedAt,
			&resolvedBy,
			&resolvedAt,
			&resolvedByUsername,
		); err != nil {
			return nil, 0, err
		}

		if resolvedBy.Valid {
			uid := domain.UserID(resolvedBy.Int32) //nolint:gosec //it's ok here
			is.ResolvedBy = &uid
		}
		if resolvedAt.Valid {
			is.ResolvedAt = &resolvedAt.Time
		}
		if resolvedByUsername.Valid {
			username := resolvedByUsername.String
			is.ResolvedByUsername = &username
		}

		issues = append(issues, is)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return issues, total, nil
}

//nolint:lll // it's ok for a big query
func (r *Repository) RecentIssues(ctx context.Context, limit uint) ([]domain.IssueExtended, error) {
	executor := r.getExecutor(ctx)
	uid := wardencontext.UserID(ctx)

	if limit == 0 {
		limit = 3
	}

	const query = `
WITH recent AS (
    SELECT
        issues.id,
        project_id,
        projects.name AS project_name,
        fingerprint,
        source,
        status,
        title,
        level,
        platform,
        first_seen,
        last_seen,
        total_events,
        issues.created_at,
        issues.updated_at
    FROM issues
    LEFT JOIN projects ON projects.id = issues.project_id
    LEFT JOIN team_members tm ON tm.team_id = projects.team_id AND tm.user_id = $1
    WHERE status = 'unresolved' AND (tm.user_id IS NOT NULL OR projects.team_id IS NULL) AND projects.archived_at IS NULL
    ORDER BY last_seen DESC
    LIMIT 100
)
SELECT *
FROM recent
ORDER BY
    /* 1-й критерий: приоритет уровня */
    CASE level
        WHEN 'fatal'     THEN 0   -- самый высокий
        WHEN 'exception' THEN 1
        WHEN 'error'     THEN 2
        WHEN 'warning'   THEN 3
        WHEN 'info'      THEN 4
        WHEN 'debug'     THEN 5
        ELSE                  6
        END,
    /* 2-й критерий: количество событий */
    total_events DESC,
    /* 3-й критерий: когда видели в последний раз */
    last_seen    DESC
LIMIT $2;`

	rows, err := executor.Query(ctx, query, uid, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	issues := make([]domain.IssueExtended, 0, limit)
	for rows.Next() {
		var is domain.IssueExtended
		if err := rows.Scan(
			&is.ID,
			&is.ProjectID,
			&is.ProjectName,
			&is.Fingerprint,
			&is.Source,
			&is.Status,
			&is.Title,
			&is.Level,
			&is.Platform,
			&is.FirstSeen,
			&is.LastSeen,
			&is.TotalEvents,
			&is.CreatedAt,
			&is.UpdatedAt,
		); err != nil {
			return nil, err
		}
		issues = append(issues, is)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return issues, nil
}

// Timeseries returns an aggregated number of issues per time bucket.
//
//nolint:gocyclo // need refactoring
func (r *Repository) Timeseries(
	ctx context.Context,
	filter *domain.IssueTimeseriesFilter,
) ([]domain.Timeseries, error) {
	executor := r.getExecutor(ctx)

	currUserID := wardencontext.UserID(ctx)
	isSuperUser := wardencontext.IsSuper(ctx)

	// ────────────────────────────────────────────────
	// 1. Calculate helper values
	// ────────────────────────────────────────────────
	if filter.Period.Interval <= 0 || filter.Period.Granularity <= 0 {
		return nil, errors.New("invalid period settings")
	}

	now := time.Now().UTC()
	from := now.Add(-filter.Period.Interval)
	bucketSec := int64(filter.Period.Granularity.Seconds())
	bucketsCnt := int(filter.Period.Interval / filter.Period.Granularity)

	// ────────────────────────────────────────────────
	// 2. Build SELECT statement dynamically
	// ────────────────────────────────────────────────
	statementBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	// `bucket` – начало интервала (rounding down)
	bucketExpr := fmt.Sprintf(
		"to_timestamp(floor(extract(epoch from last_seen) / %d) * %d)::timestamptz AS bucket",
		bucketSec, bucketSec,
	)

	selectCols := []string{
		bucketExpr,
		"COUNT(*) AS occurrences",
	}

	// Optional grouping column + human-readable alias `name`
	groupColumn := ""
	switch filter.GroupBy {
	case domain.IssueTimeseriesGroupProject:
		groupColumn = "project_id"
	case domain.IssueTimeseriesGroupIssue:
		groupColumn = "id"
	case domain.IssueTimeseriesGroupLevel:
		groupColumn = "level"
	case domain.IssueTimeseriesGroupStatus:
		groupColumn = "status"
	default:
		// no extra column
	}

	if groupColumn != "" {
		selectCols = append(selectCols, groupColumn+" AS name")
	}

	qb := statementBuilder.
		Select(selectCols...).
		From("issues").
		LeftJoin("projects ON issues.project_id = projects.id").
		Where(sq.GtOrEq{"last_seen": from}).
		Where(sq.LtOrEq{"last_seen": now}).
		Where("projects.archived_at IS NULL")

	// Filters ---------------------------------------------------------------
	if filter.ProjectID != nil {
		qb = qb.Where(sq.Eq{"project_id": *filter.ProjectID})
	}
	if filter.IssueID != nil {
		qb = qb.Where(sq.Eq{"id": *filter.IssueID})
	}
	if len(filter.Levels) > 0 {
		qb = qb.Where(sq.Eq{"level": filter.Levels})
	}
	if len(filter.Statuses) > 0 {
		qb = qb.Where(sq.Eq{"status": filter.Statuses})
	}
	if !isSuperUser {
		qb = qb.LeftJoin(
			"team_members tm ON tm.team_id = projects.team_id AND tm.user_id = ?", currUserID)
		qb = qb.Where("tm.user_id IS NOT NULL OR projects.team_id IS NULL")
	}

	// Group BY --------------------------------------------------------------
	groupBy := []string{"bucket"}
	if groupColumn != "" {
		groupBy = append(groupBy, groupColumn)
	}
	qb = qb.GroupBy(groupBy...).OrderBy("bucket")

	query, args, err := qb.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build timeseries query: %w", err)
	}

	// ────────────────────────────────────────────────
	// 3. Execute and collect
	// ────────────────────────────────────────────────
	rows, err := executor.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("timeseries query: %w", err)
	}
	defer rows.Close()

	// Map[name]→counts[]
	type key = string
	acc := make(map[key][]uint)

	for rows.Next() {
		var (
			bucket      time.Time
			occurrences uint
			name        string
		)

		if groupColumn == "" {
			// when no grouping, we still scan into dummy `name`
			name = "all"
			if err := rows.Scan(&bucket, &occurrences); err != nil {
				return nil, fmt.Errorf("scan row: %w", err)
			}
		} else {
			if err := rows.Scan(&bucket, &occurrences, &name); err != nil {
				return nil, fmt.Errorf("scan row: %w", err)
			}
		}

		// initialise slice for the group if needed
		if _, ok := acc[name]; !ok {
			acc[name] = make([]uint, bucketsCnt)
		}

		// Calculate index of the bucket
		idx := int(bucket.Sub(from) / filter.Period.Granularity)
		if idx >= 0 && idx < bucketsCnt {
			acc[name][idx] = occurrences
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rows: %w", err)
	}

	// ────────────────────────────────────────────────
	// 4. Convert to domain objects
	// ────────────────────────────────────────────────
	result := make([]domain.Timeseries, 0, len(acc))
	for name, counts := range acc {
		result = append(result, domain.Timeseries{
			Name:        name,
			Period:      filter.Period,
			Occurrences: counts,
		})
	}

	// deterministic order (useful for tests / clients)
	slices.SortFunc(result, func(a, b domain.Timeseries) int {
		return strings.Compare(a.Name, b.Name)
	})

	return result, nil
}

func (r *Repository) CountForLevels(
	ctx context.Context,
	projectID domain.ProjectID,
	period time.Duration,
) (map[domain.IssueLevel]uint64, error) {
	executor := r.getExecutor(ctx)

	stmt := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Select("level", "COUNT(*)").
		From("issues").
		Join("projects ON issues.project_id = projects.id").
		Where("issues.project_id = ?", projectID).
		Where("last_seen >= ?", time.Now().Add(-period)).
		Where("projects.archived_at IS NULL").
		GroupBy("level")

	sqlStr, args, err := stmt.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build sql: %w", err)
	}

	res := map[domain.IssueLevel]uint64{
		domain.IssueLevelFatal:     0,
		domain.IssueLevelException: 0,
		domain.IssueLevelError:     0,
		domain.IssueLevelWarning:   0,
		domain.IssueLevelInfo:      0,
		domain.IssueLevelDebug:     0,
	}

	rows, err := executor.Query(ctx, sqlStr, args...)
	if err != nil {
		return nil, fmt.Errorf("query count for levels: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var lvl string
		var cnt uint64
		if err := rows.Scan(&lvl, &cnt); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}

		switch domain.IssueLevel(lvl) {
		case domain.IssueLevelException,
			domain.IssueLevelFatal,
			domain.IssueLevelDebug,
			domain.IssueLevelError,
			domain.IssueLevelWarning,
			domain.IssueLevelInfo:
			res[domain.IssueLevel(lvl)] = cnt
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rows: %w", err)
	}

	return res, nil
}

func (r *Repository) MostFrequent(
	ctx context.Context,
	projectID domain.ProjectID,
	period time.Duration,
	limit uint,
) ([]domain.IssueExtended, error) {
	executor := r.getExecutor(ctx)

	if limit == 0 {
		limit = 5
	}

	stmt := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Select(
			"issues.id AS id",
			"project_id",
			"projects.name AS project_name",
			"fingerprint",
			"source",
			"status",
			"title",
			"level",
			"platform",
			"first_seen",
			"last_seen",
			"total_events",
			"issues.created_at AS created_at",
			"issues.updated_at AS updated_at",
		).
		From("issues").
		LeftJoin("projects ON issues.project_id = projects.id").
		Where("project_id = ?", projectID).
		Where("last_seen >= ?", time.Now().Add(-period)).
		Where("projects.archived_at IS NULL").
		OrderBy(
			`CASE level 
				WHEN 'fatal' THEN 1 
				WHEN 'exception' THEN 2 
				WHEN 'error' THEN 3 
				WHEN 'warning' THEN 4 
				WHEN 'info' THEN 5 
				WHEN 'debug' THEN 6 
				ELSE 7 
			END`,
			"total_events DESC",
			"last_seen DESC",
		).
		Limit(uint64(limit))

	sqlStr, args, err := stmt.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build sql: %w", err)
	}

	rows, err := executor.Query(ctx, sqlStr, args...)
	if err != nil {
		return nil, fmt.Errorf("query most frequent: %w", err)
	}
	defer rows.Close()

	issues := make([]domain.IssueExtended, 0, limit)
	for rows.Next() {
		var is domain.IssueExtended
		if err := rows.Scan(
			&is.ID,
			&is.ProjectID,
			&is.ProjectName,
			&is.Fingerprint,
			&is.Source,
			&is.Status,
			&is.Title,
			&is.Level,
			&is.Platform,
			&is.FirstSeen,
			&is.LastSeen,
			&is.TotalEvents,
			&is.CreatedAt,
			&is.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		issues = append(issues, is)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rows: %w", err)
	}

	return issues, nil
}

func (r *Repository) UpdateStatus(ctx context.Context, issueID domain.IssueID, status domain.IssueStatus) error {
	executor := r.getExecutor(ctx)
	query := "UPDATE issues SET status = $1, updated_at = NOW() WHERE id = $2"

	_, err := executor.Exec(ctx, query, status, issueID)
	if err != nil {
		return fmt.Errorf("exec update: %w", err)
	}

	return nil
}

func (r *Repository) MarkAsNotified(ctx context.Context, issueID domain.IssueID) error {
	executor := r.getExecutor(ctx)
	const query = "UPDATE issues SET last_notification_at = NOW(), updated_at = NOW() WHERE id = $1"

	_, err := executor.Exec(ctx, query, issueID)
	if err != nil {
		return fmt.Errorf("exec update: %w", err)
	}

	return nil
}

func (r *Repository) ListUnresolved(ctx context.Context) ([]domain.IssueExtended, error) {
	executor := r.getExecutor(ctx)
	const query = `
SELECT 
    i.*,
    p.name AS project_name,
    res.resolved_at,
    res.resolved_by,
    u.username AS resolved_by_username
FROM issues i
JOIN projects p ON p.id = i.project_id
JOIN LATERAL (
    SELECT id
    FROM notification_settings ns
    WHERE ns.project_id = p.id AND ns.enabled = true
    LIMIT 1
    ) ns ON true
LEFT JOIN LATERAL (
    SELECT resolved_at, resolved_by
    FROM resolutions
    WHERE resolutions.issue_id = i.id
    ORDER BY resolved_at DESC
    LIMIT 1
    ) res ON true
LEFT JOIN users u ON u.id = res.resolved_by
JOIN LATERAL (
    SELECT id
    FROM notification_rules
    WHERE notification_rules.notification_setting_id = ns.id AND i.level = notification_rules.event_level
    LIMIT 1
    ) rul on true
WHERE
    i.status = 'unresolved'
ORDER BY i.id`
	rows, err := executor.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	list, err := pgx.CollectRows(rows, pgx.RowToStructByName[issueExtendedModel])
	if err != nil {
		return nil, fmt.Errorf("collect rows: %w", err)
	}

	result := make([]domain.IssueExtended, 0, len(list))
	for i := range list {
		elem := list[i]
		issueExtended := elem.toDomain()
		result = append(result, issueExtended)
	}

	return result, nil
}

// NewIssuesForRelease returns issues first seen in this release.
func (r *Repository) NewIssuesForRelease(
	ctx context.Context,
	projectID domain.ProjectID,
	release string,
) ([]string, error) {
	executor := r.getExecutor(ctx)
	const query = `
SELECT ir.issue_id
FROM issue_releases ir
JOIN issues i ON i.id = ir.issue_id
JOIN releases r ON r.id = ir.release_id
WHERE i.project_id = $1 AND r.version = $2 AND ir.first_seen_in = true
`
	rows, err := executor.Query(ctx, query, projectID, release)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []string
	for rows.Next() {
		var issueID uint
		if err := rows.Scan(&issueID); err != nil {
			return nil, err
		}
		result = append(result, strconv.FormatUint(uint64(issueID), 10))
	}

	return result, nil
}

// ResolvedInRelease returns issues resolved in this release.
func (r *Repository) ResolvedInRelease(
	ctx context.Context,
	projectID domain.ProjectID,
	release string,
) ([]string, error) {
	executor := r.getExecutor(ctx)
	const query = `
SELECT ir.issue_id
FROM issue_releases ir
JOIN issues i ON i.id = ir.issue_id
JOIN releases r ON r.id = ir.release_id
WHERE i.project_id = $1 AND r.version = $2 AND i.status = 'resolved'
`
	rows, err := executor.Query(ctx, query, projectID, release)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []string
	for rows.Next() {
		var issueID uint
		if err := rows.Scan(&issueID); err != nil {
			return nil, err
		}
		result = append(result, strconv.FormatUint(uint64(issueID), 10))
	}

	return result, nil
}

// FixTimesForRelease returns fix times for issues resolved in this release.
func (r *Repository) FixTimesForRelease(
	ctx context.Context,
	projectID domain.ProjectID,
	release string,
) (map[string]time.Duration, error) {
	executor := r.getExecutor(ctx)
	const query = `
SELECT ir.issue_id, EXTRACT(EPOCH FROM (i.updated_at - i.first_seen))
FROM issue_releases ir
JOIN issues i ON i.id = ir.issue_id
JOIN releases r ON r.id = ir.release_id
WHERE i.project_id = $1 AND r.version = $2 AND i.status = 'resolved'
`
	rows, err := executor.Query(ctx, query, projectID, release)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make(map[string]time.Duration)
	for rows.Next() {
		var issueID uint
		var seconds float64
		if err := rows.Scan(&issueID, &seconds); err != nil {
			return nil, err
		}
		result[strconv.FormatUint(uint64(issueID), 10)] = time.Duration(seconds) * time.Second
	}

	return result, nil
}

// RegressionsForRelease returns issues that were resolved but re-appeared in this release.
func (r *Repository) RegressionsForRelease(
	ctx context.Context,
	projectID domain.ProjectID,
	release string,
) ([]string, error) {
	executor := r.getExecutor(ctx)
	const query = `
SELECT DISTINCT ir.issue_id
FROM issue_releases ir
JOIN issues i ON i.id = ir.issue_id
JOIN releases r ON r.id = ir.release_id
LEFT JOIN resolutions res ON res.issue_id = i.id
WHERE i.project_id = $1
  AND r.version = $2
  AND res.resolved_at IS NOT NULL
  AND i.last_seen > res.resolved_at
  AND i.status != 'resolved'
`
	rows, err := executor.Query(ctx, query, projectID, release)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []string
	for rows.Next() {
		var issueID uint
		if err := rows.Scan(&issueID); err != nil {
			return nil, err
		}
		result = append(result, strconv.FormatUint(uint64(issueID), 10))
	}

	return result, nil
}

//nolint:ireturn // it's ok here
func (r *Repository) getExecutor(ctx context.Context) db.Tx {
	if tx := db.TxFromContext(ctx); tx != nil {
		return tx
	}

	return r.db
}
