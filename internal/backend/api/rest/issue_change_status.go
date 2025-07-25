package rest

import (
	"context"
	"errors"
	"log/slog"

	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

func (r *RestAPI) ChangeIssueStatus(
	ctx context.Context,
	req *generatedapi.ChangeIssueStatusReq,
	params generatedapi.ChangeIssueStatusParams,
) (generatedapi.ChangeIssueStatusRes, error) {
	issueID := domain.IssueID(params.IssueID)

	// Check if the user has permission to manage the issue
	if err := r.permissionsService.CanManageIssue(ctx, issueID); err != nil {
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

	err := r.issueUseCase.ChangeStatus(ctx, issueID, domain.IssueStatus(req.Status))
	if err != nil {
		slog.Error("change issue status failed", "error", err)

		switch {
		case errors.Is(err, domain.ErrEntityNotFound):
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString(err.Error()),
			}}, nil
		default:
			return nil, err
		}
	}

	return &generatedapi.ChangeIssueStatusNoContent{}, nil
}
