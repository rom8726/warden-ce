package rest

import (
	"context"
	"log/slog"

	"github.com/rom8726/warden/internal/backend/dto"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

func (r *RestAPI) GetRecentIssues(
	ctx context.Context,
	params generatedapi.GetRecentIssuesParams,
) (generatedapi.GetRecentIssuesRes, error) {
	list, err := r.issueUseCase.RecentIssues(ctx, params.Limit)
	if err != nil {
		slog.Error("get recent issues failed", "error", err)

		return nil, err
	}

	items := make([]generatedapi.IssueSummary, 0, len(list))
	for i := range list {
		elem := list[i]
		items = append(items, generatedapi.IssueSummary{
			ID:        elem.ID.Uint(),
			ProjectID: elem.ProjectID.Uint(),
			Title:     elem.Title,
			Level:     dto.DomainLevelToAPI(elem.Level),
			Count:     elem.TotalEvents,
			LastSeen:  elem.LastSeen,
		})
	}

	return &generatedapi.ListIssueSummariesResponse{Issues: items}, nil
}
