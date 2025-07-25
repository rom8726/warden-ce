package dto

import (
	"time"

	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

func ToReleaseAnalyticsSummary(rel domain.Release, stats domain.ReleaseStats) generatedapi.ReleaseAnalyticsSummary {
	return generatedapi.ReleaseAnalyticsSummary{
		Version:                rel.Version,
		CreatedAt:              rel.CreatedAt,
		KnownIssuesTotal:       stats.KnownIssuesTotal,
		NewIssuesTotal:         stats.NewIssuesTotal,
		RegressionsTotal:       stats.RegressionsTotal,
		ResolvedInVersionTotal: stats.ResolvedInVersionTotal,
		UsersAffected:          stats.UsersAffected,
	}
}

func ToReleaseAnalyticsDetails(details domain.ReleaseAnalyticsDetails) generatedapi.ReleaseAnalyticsDetails {
	return generatedapi.ReleaseAnalyticsDetails{
		Version:              details.Release.Version,
		CreatedAt:            details.Release.CreatedAt,
		Stats:                ToReleaseAnalyticsSummary(details.Release, details.Stats),
		TopIssues:            ToIssueSummaries(details.TopIssues),
		SeverityDistribution: details.Stats.SeverityDistribution,
		FixTime: generatedapi.ReleaseAnalyticsDetailsFixTime{
			Avg:    toOptFloat32(details.Stats.AvgFixTime),
			Median: toOptFloat32(details.Stats.MedianFixTime),
			P95:    toOptFloat32(details.Stats.P95FixTime),
		},
		Segments: generatedapi.ReleaseAnalyticsDetailsSegments{
			Platform:    generatedapi.NewOptReleaseAnalyticsDetailsSegmentsPlatform(details.ByPlatform),
			BrowserName: generatedapi.NewOptReleaseAnalyticsDetailsSegmentsBrowserName(details.ByBrowser),
			OsName:      generatedapi.NewOptReleaseAnalyticsDetailsSegmentsOsName(details.ByOS),
			DeviceArch:  generatedapi.NewOptReleaseAnalyticsDetailsSegmentsDeviceArch(details.ByDeviceArch),
			RuntimeName: generatedapi.NewOptReleaseAnalyticsDetailsSegmentsRuntimeName(details.ByRuntimeName),
		},
	}
}

func ToReleaseComparison(comp domain.ReleaseComparison) generatedapi.ReleaseComparison {
	return generatedapi.ReleaseComparison{
		Base:   ToReleaseAnalyticsSummary(domain.Release{Version: comp.BaseRelease.Release}, comp.BaseRelease),
		Target: ToReleaseAnalyticsSummary(domain.Release{Version: comp.TargetRelease.Release}, comp.TargetRelease),
		Delta: generatedapi.ReleaseComparisonDelta{
			KnownIssuesTotal:       toOptUInt(comp.Delta["known_issues_total"]),
			NewIssuesTotal:         toOptUInt(comp.Delta["new_issues_total"]),
			RegressionsTotal:       toOptUInt(comp.Delta["regressions_total"]),
			ResolvedInVersionTotal: toOptUInt(comp.Delta["resolved_in_version_total"]),
			UsersAffected:          toOptUInt(comp.Delta["users_affected"]),
		},
	}
}

func ToReleaseSegmentsResponse(segment string, values map[string]uint) generatedapi.ReleaseSegmentsResponse {
	return generatedapi.ReleaseSegmentsResponse{
		Segment: segment,
		Values:  values,
	}
}

func toOptFloat32(val *time.Duration) generatedapi.OptFloat32 {
	if val == nil {
		return generatedapi.OptFloat32{}
	}

	return generatedapi.NewOptFloat32(float32(val.Seconds()))
}

func toOptUInt(val uint) generatedapi.OptUint {
	return generatedapi.NewOptUint(val)
}

// ToIssueSummaries — конвертация []domain.IssueSummary (или []string) в []generatedapi.IssueSummary.
func ToIssueSummaries(issues []domain.Issue) []generatedapi.IssueSummary {
	res := make([]generatedapi.IssueSummary, len(issues))
	for i, issue := range issues {
		res[i] = generatedapi.IssueSummary{
			ID:        issue.ID.Uint(),
			ProjectID: issue.ProjectID.Uint(),
			Title:     issue.Title,
			Level:     generatedapi.IssueLevel(issue.Level),
			Count:     issue.TotalEvents,
			LastSeen:  issue.LastSeen,
		}
	}

	return res
}
