package rest

import (
	"context"
	"errors"
	"log/slog"

	dto2 "github.com/rom8726/warden/internal/backend/dto"
	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

func (r *RestAPI) GetProjectStats(
	ctx context.Context,
	params generatedapi.GetProjectStatsParams,
) (generatedapi.GetProjectStatsRes, error) {
	period, err := dto2.ParseHumanDuration(string(params.Period))
	if err != nil {
		slog.Error("invalid period", "error", err)

		return nil, err
	}

	stats, err := r.projectsUseCase.GeneralStats(ctx, domain.ProjectID(params.ProjectID), period)
	if err != nil {
		slog.Error("get project stats failed", "error", err)

		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString(err.Error()),
			}}, nil
		}

		return nil, err
	}

	issues := make([]generatedapi.IssueSummary, 0, len(stats.MostFrequentIssues))
	for i := range stats.MostFrequentIssues {
		elem := stats.MostFrequentIssues[i]
		issues = append(issues, generatedapi.IssueSummary{
			ID:        elem.ID.Uint(),
			ProjectID: elem.ProjectID.Uint(),
			Title:     elem.Title,
			Level:     dto2.DomainLevelToAPI(elem.Level),
			Count:     elem.TotalEvents,
			LastSeen:  elem.LastSeen,
		})
	}

	return &generatedapi.ProjectStatsResponse{
		TotalIssues: stats.TotalIssues,
		IssuesByLevel: generatedapi.ProjectStatsResponseIssuesByLevel{
			Fatal:     stats.FatalIssues,
			Exception: stats.ExceptionIssues,
			Error:     stats.ErrorIssues,
			Warning:   stats.WarningIssues,
			Info:      stats.InfoIssues,
			Debug:     stats.DebugIssues,
		},
		MostFrequentIssues: issues,
	}, nil
}
