package releasestats

import (
	"encoding/json"
	"time"

	"github.com/rom8726/warden/internal/domain"
)

type releaseStatsModel struct {
	ProjectID   uint      `db:"project_id"`
	ReleaseID   uint      `db:"release_id"`
	Release     string    `db:"release"`
	GeneratedAt time.Time `db:"generated_at"`

	KnownIssuesTotal uint `db:"known_issues_total"`
	NewIssuesTotal   uint `db:"new_issues_total"`
	RegressionsTotal uint `db:"regressions_total"`

	ResolvedInVersionTotal uint `db:"resolved_in_version_total"`
	FixedNewInVersionTotal uint `db:"fixed_new_in_version_total"`
	FixedOldInVersionTotal uint `db:"fixed_old_in_version_total"`

	AvgFixTime    string `db:"avg_fix_time"`
	MedianFixTime string `db:"median_fix_time"`
	P95FixTime    string `db:"p95_fix_time"`

	SeverityDistribution []byte `db:"severity_distribution"`
	UsersAffected        uint   `db:"users_affected"`
}

func (m *releaseStatsModel) toDomain() (domain.ReleaseStats, error) {
	avgFixRef, err := parseStringToDurationPtr(m.AvgFixTime)
	if err != nil {
		return domain.ReleaseStats{}, err
	}
	medianFixRef, err := parseStringToDurationPtr(m.MedianFixTime)
	if err != nil {
		return domain.ReleaseStats{}, err
	}
	p95FixRef, err := parseStringToDurationPtr(m.P95FixTime)
	if err != nil {
		return domain.ReleaseStats{}, err
	}

	severity := make(map[string]uint)
	_ = json.Unmarshal(m.SeverityDistribution, &severity)

	return domain.ReleaseStats{
		ProjectID:              domain.ProjectID(m.ProjectID),
		ReleaseID:              domain.ReleaseID(m.ReleaseID),
		Release:                m.Release,
		GeneratedAt:            m.GeneratedAt,
		KnownIssuesTotal:       m.KnownIssuesTotal,
		NewIssuesTotal:         m.NewIssuesTotal,
		RegressionsTotal:       m.RegressionsTotal,
		ResolvedInVersionTotal: m.ResolvedInVersionTotal,
		FixedNewInVersionTotal: m.FixedNewInVersionTotal,
		FixedOldInVersionTotal: m.FixedOldInVersionTotal,
		AvgFixTime:             avgFixRef,
		MedianFixTime:          medianFixRef,
		P95FixTime:             p95FixRef,
		SeverityDistribution:   severity,
		UsersAffected:          m.UsersAffected,
	}, nil
}

//nolint:nilnil //it's ok here
func parseStringToDurationPtr(str string) (*time.Duration, error) {
	if str == "" || str == "0s" {
		return nil, nil
	}

	duration, err := time.ParseDuration(str)
	if err != nil {
		return nil, err
	}

	if duration == 0 {
		return nil, nil
	}

	return &duration, nil
}
