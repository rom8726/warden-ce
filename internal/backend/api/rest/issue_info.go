package rest

import (
	"context"
	"errors"
	"log/slog"

	"github.com/rom8726/warden/internal/backend/dto"
	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

func (r *RestAPI) GetIssue(ctx context.Context, params generatedapi.GetIssueParams) (generatedapi.GetIssueRes, error) {
	issueID := domain.IssueID(params.IssueID)

	// Check if the user has access to the issue
	if err := r.permissionsService.CanAccessIssue(ctx, issueID); err != nil {
		slog.Error("permission denied", "error", err, "issue_id", issueID)

		if errors.Is(err, domain.ErrPermissionDenied) {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString("issue not found"),
			}}, nil
		}

		if errors.Is(err, domain.ErrUserNotFound) {
			return &generatedapi.ErrorUnauthorized{Error: generatedapi.ErrorUnauthorizedError{
				Message: generatedapi.NewOptString("unauthorized"),
			}}, nil
		}

		return nil, err
	}

	issue, err := r.issueUseCase.GetByIDWithChildren(ctx, issueID)
	if err != nil {
		slog.Error("get issue failed", "error", err)

		return nil, err
	}

	resp := dto.MakeIssueResponseWithEvent(issue)

	return &resp, nil
}
