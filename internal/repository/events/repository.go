//nolint:errcheck // for clickhouse rows
package events

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/rom8726/warden/internal/domain"
	"github.com/rom8726/warden/internal/infra"
	"github.com/rom8726/warden/pkg/kafka"
)

type Repository struct {
	clickHouseClient infra.ClickHouseConn
	producer         kafka.DataProducer
}

func New(clickHouseClient infra.ClickHouseConn, producer *kafka.TopicProducer) *Repository {
	return &Repository{
		clickHouseClient: clickHouseClient,
		producer:         producer,
	}
}

func (r *Repository) Store(ctx context.Context, event *domain.Event) error {
	model, err := fromDomain(event)
	if err != nil {
		return fmt.Errorf("convert domain event to repo model: %w", err)
	}

	data, err := json.Marshal(model)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	if err := r.producer.Produce(ctx, data); err != nil {
		return fmt.Errorf("produce event: %w", err)
	}

	return nil
}

// StoreWithFingerprints calculates fingerprints for the event and stores it.
func (r *Repository) StoreWithFingerprints(ctx context.Context, event *domain.Event) error {
	// Ensure we have a timestamp
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	if event.GroupHash == "" {
		event.GroupHash = event.FullFingerprint()
	}

	return r.Store(ctx, event)
}

func (r *Repository) FetchForIssue(
	ctx context.Context,
	projectID domain.ProjectID,
	fingerprint string,
	limit uint,
) ([]domain.Event, error) {
	query := `
SELECT *
FROM events
WHERE project_id = ? AND group_hash = ?
ORDER BY timestamp DESC
LIMIT ?`
	rows, err := r.clickHouseClient.QueryWithRetries(ctx, query, projectID, fingerprint, limit)
	if err != nil {
		return nil, fmt.Errorf("query events: %w", err)
	}
	defer rows.Close()

	var events []domain.Event
	for rows.Next() {
		var evModel eventModel
		if err := rows.ScanStruct(&evModel); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}

		event, err := evModel.toDomain()
		if err != nil {
			return nil, fmt.Errorf("convert event model to domain: %w", err)
		}
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rows: %w", err)
	}

	return events, nil
}

//nolint:gocyclo // need refactoring
func (r *Repository) Timeseries(
	ctx context.Context,
	filter *domain.EventTimeseriesFilter,
) ([]domain.Timeseries, error) {
	if filter == nil {
		return nil, errors.New("filter is nil")
	}
	gran := filter.Period.Granularity
	if gran == 0 {
		return nil, errors.New("granularity is zero")
	}

	// 1. Вычисляем границы периода
	nowAligned := time.Now().Truncate(time.Minute)
	start := nowAligned.Add(-filter.Period.Interval)
	end := nowAligned

	// 2. Вспомогательные выражения
	granExpr, err := granularityExpr(gran)
	if err != nil {
		return nil, err
	}
	groupExpr, groupName := groupExpression(filter.GroupBy)

	// 3. Строим SQL при помощи Squirrel
	sb := sq.StatementBuilder.PlaceholderFormat(sq.Question) // builder

	qb := sb.
		Select(
			groupExpr+" AS group_key",
			granExpr+" AS bucket",
			"count() AS cnt",
		).
		From("events").
		Where(sq.GtOrEq{"timestamp": start}).
		Where(sq.Lt{"timestamp": end}).
		GroupBy("group_key", "bucket").
		OrderBy("group_key", "bucket")

	if filter.ProjectID != nil {
		qb = qb.Where(sq.Eq{"project_id": *filter.ProjectID})
	}
	if len(filter.Levels) > 0 {
		qb = qb.Where(sq.Eq{"level": filter.Levels})
	}
	if filter.Release != nil {
		qb = qb.Where(sq.Eq{"release": *filter.Release})
	}

	query, args, err := qb.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build timeseries query: %w", err)
	}

	// 4. Выполняем запрос
	rows, err := r.clickHouseClient.QueryWithRetries(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query timeseries: %w", err)
	}
	defer rows.Close()

	bucketsTotal := int(filter.Period.Interval / gran)
	seriesMap := make(map[string][]uint, 8)

	for rows.Next() {
		var (
			name   string
			bucket time.Time
			count  uint64
		)
		if err := rows.Scan(&name, &bucket, &count); err != nil {
			return nil, fmt.Errorf("scan timeseries row: %w", err)
		}

		idx := int(bucket.Sub(start) / gran)
		if idx < 0 || idx >= bucketsTotal {
			continue
		}

		if _, ok := seriesMap[name]; !ok {
			seriesMap[name] = make([]uint, bucketsTotal)
		}
		seriesMap[name][idx] = uint(count)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rows: %w", err)
	}

	// 5. Преобразуем в доменную модель
	result := make([]domain.Timeseries, 0, len(seriesMap))
	for name, points := range seriesMap {
		result = append(result, domain.Timeseries{
			Name:        name,
			Period:      filter.Period,
			Occurrences: points,
		})
	}
	if filter.GroupBy == domain.EventTimeseriesGroupNone && len(result) == 0 {
		result = append(result, domain.Timeseries{
			Name:        groupName,
			Period:      filter.Period,
			Occurrences: make([]uint, bucketsTotal),
		})
	}

	return result, nil
}

//nolint:gocyclo // need refactoring
func (r *Repository) IssueTimeseries(
	ctx context.Context,
	fingerprint string,
	filter *domain.IssueEventsTimeseriesFilter,
) ([]domain.Timeseries, error) {
	// --- 0. Валидация -------------------------------------------------------
	if fingerprint == "" {
		return nil, errors.New("fingerprint is empty")
	}
	if filter == nil {
		return nil, errors.New("filter is nil")
	}
	if filter.Period.Granularity == 0 {
		return nil, errors.New("granularity is zero")
	}

	gran := filter.Period.Granularity

	// --- 1. Период запроса --------------------------------------------------
	nowAligned := time.Now().Truncate(time.Minute)
	start := nowAligned.Add(-filter.Period.Interval)
	end := nowAligned

	// --- 2. Выражения ClickHouse -------------------------------------------
	granExpr, err := granularityExpr(gran)
	if err != nil {
		return nil, err
	}
	groupExpr, groupName := groupExpression(filter.GroupBy)

	// --- 3. QueryBuilder ----------------------------------------------------
	sb := sq.StatementBuilder.PlaceholderFormat(sq.Question)

	qb := sb.
		Select(
			groupExpr+" AS group_key",
			granExpr+" AS bucket",
			"count() AS cnt",
		).
		From("events").
		Where(sq.Eq{"group_hash": fingerprint}).
		Where(sq.GtOrEq{"timestamp": start}).
		Where(sq.Lt{"timestamp": end}).
		GroupBy("group_key", "bucket").
		OrderBy("group_key", "bucket")

	if filter.ProjectID != 0 {
		qb = qb.Where(sq.Eq{"project_id": filter.ProjectID})
	}
	if len(filter.Levels) > 0 {
		qb = qb.Where(sq.Eq{"level": filter.Levels})
	}

	query, args, err := qb.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build issue-timeseries query: %w", err)
	}

	// --- 4. Выполнение ------------------------------------------------------
	rows, err := r.clickHouseClient.QueryWithRetries(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query issue-timeseries: %w", err)
	}
	defer rows.Close()

	bucketsTotal := int(filter.Period.Interval / gran)
	seriesMap := make(map[string][]uint, 4)

	for rows.Next() {
		var (
			name   string
			bucket time.Time
			cnt    uint64
		)
		if err := rows.Scan(&name, &bucket, &cnt); err != nil {
			return nil, fmt.Errorf("scan issue-timeseries row: %w", err)
		}

		idx := int(bucket.Sub(start) / gran)
		if idx < 0 || idx >= bucketsTotal {
			continue
		}

		if _, ok := seriesMap[name]; !ok {
			seriesMap[name] = make([]uint, bucketsTotal)
		}
		seriesMap[name][idx] = uint(cnt)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rows: %w", err)
	}

	// --- 5. В доменную модель ----------------------------------------------
	result := make([]domain.Timeseries, 0, len(seriesMap))
	for name, points := range seriesMap {
		result = append(result, domain.Timeseries{
			Name:        name,
			Period:      filter.Period,
			Occurrences: points,
		})
	}

	// Если группировки нет, но результатов нет – вернуть пустой ряд,
	// чтобы фронт всегда получал ожидаемое количество bucket'ов.
	if filter.GroupBy == domain.EventTimeseriesGroupNone && len(result) == 0 {
		result = append(result, domain.Timeseries{
			Name:        groupName,
			Period:      filter.Period,
			Occurrences: make([]uint, bucketsTotal),
		})
	}

	return result, nil
}

func granularityExpr(gran time.Duration) (string, error) {
	if gran <= 0 {
		return "", errors.New("granularity must be positive")
	}

	const day = 24 * time.Hour

	switch {
	// Кратность дню
	case gran%day == 0:
		n := int64(gran / day)

		return fmt.Sprintf("toStartOfInterval(timestamp, INTERVAL %d DAY)", n), nil

	// Кратность часу
	case gran%time.Hour == 0:
		n := int64(gran / time.Hour)

		return fmt.Sprintf("toStartOfInterval(timestamp, INTERVAL %d HOUR)", n), nil

	// Кратность минуте
	case gran%time.Minute == 0:
		n := int64(gran / time.Minute)

		return fmt.Sprintf("toStartOfInterval(timestamp, INTERVAL %d MINUTE)", n), nil
	}

	return "", fmt.Errorf("unsupported granularity %s (must be multiple of 1m,1h or 1d)", gran)
}

func groupExpression(g domain.EventTimeseriesGroup) (expr, defaultName string) {
	switch g {
	case domain.EventTimeseriesGroupProject:
		return "toString(project_id)", ""
	case domain.EventTimeseriesGroupLevel:
		return "level", ""
	default:
		// None
		return "''", "all"
	}
}

func (r *Repository) EventsByRelease(
	ctx context.Context,
	projectID domain.ProjectID,
	release string,
	limit uint,
) ([]domain.Event, error) {
	const query = `
SELECT * FROM events
WHERE project_id = ? AND release = ?
ORDER BY timestamp DESC
LIMIT ?`
	rows, err := r.clickHouseClient.QueryWithRetries(ctx, query, projectID, release, limit)
	if err != nil {
		return nil, fmt.Errorf("query events by release: %w", err)
	}
	defer rows.Close()

	var events []domain.Event
	for rows.Next() {
		var evModel eventModel
		if err := rows.ScanStruct(&evModel); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		event, err := evModel.toDomain()
		if err != nil {
			return nil, fmt.Errorf("convert event model to domain: %w", err)
		}
		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rows: %w", err)
	}

	return events, nil
}

func (r *Repository) TopIssuesByRelease(
	ctx context.Context,
	projectID domain.ProjectID,
	release string,
	limit uint,
) ([]string, error) {
	const query = `
SELECT group_hash, count() as cnt
FROM events
WHERE project_id = ? AND release = ?
GROUP BY group_hash
ORDER BY cnt DESC
LIMIT ?`
	rows, err := r.clickHouseClient.QueryWithRetries(ctx, query, projectID, release, limit)
	if err != nil {
		return nil, fmt.Errorf("query top issues by release: %w", err)
	}
	defer rows.Close()

	var issues []string
	for rows.Next() {
		var groupHash string
		var cnt uint64
		if err := rows.Scan(&groupHash, &cnt); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		issues = append(issues, groupHash)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rows: %w", err)
	}

	return issues, nil
}

func (r *Repository) AggregateBySegment(
	ctx context.Context,
	projectID domain.ProjectID,
	release string,
	segment domain.SegmentName,
) (map[string]uint, error) {
	query := fmt.Sprintf(`
SELECT %s, count() as cnt
FROM events
WHERE project_id = ? AND release = ?
GROUP BY %s
ORDER BY cnt DESC`, segment, segment)
	rows, err := r.clickHouseClient.QueryWithRetries(ctx, query, projectID, release)
	if err != nil {
		return nil, fmt.Errorf("aggregate by segment: %w", err)
	}
	defer rows.Close()

	result := make(map[string]uint)
	for rows.Next() {
		var key string
		var cnt uint64
		if err := rows.Scan(&key, &cnt); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		result[key] = uint(cnt)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rows: %w", err)
	}

	return result, nil
}
