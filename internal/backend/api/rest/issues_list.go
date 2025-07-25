package rest

import (
	"context"
	"log/slog"

	"github.com/rom8726/warden/internal/backend/dto"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

func (r *RestAPI) ListIssues(
	ctx context.Context,
	params generatedapi.ListIssuesParams,
) (generatedapi.ListIssuesRes, error) {
	filter := dto.MakeIssuesListFilter(params)

	list, total, err := r.issueUseCase.List(ctx, &filter)
	if err != nil {
		slog.Error("list issues failed", "error", err)

		return nil, err
	}

	issues := make([]generatedapi.Issue, 0, len(list))
	for i := range list {
		elem := list[i]
		issues = append(issues, dto.DomainIssueToAPI(
			elem.Issue,
			elem.ProjectName,
			elem.ResolvedAt,
			elem.ResolvedBy,
			elem.ResolvedByUsername,
		))
	}

	return &generatedapi.ListIssuesResponse{
		Issues:  issues,
		Total:   uint(total),
		Page:    params.Page,
		PerPage: params.PerPage,
	}, nil
}
