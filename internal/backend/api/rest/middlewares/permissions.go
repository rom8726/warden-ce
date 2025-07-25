package middlewares

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/rom8726/warden/internal/backend/contract"
	wardencontext "github.com/rom8726/warden/internal/context"
	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

const (
	RecentKeyword = "recent"
	issuesStr     = "issues"
)

// ProjectAccess middleware checks if the user has access to the project.
func ProjectAccess(permissionsService contract.PermissionsService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if request.Method == http.MethodPost {
				next.ServeHTTP(writer, request)

				return
			}

			// Extract project ID from URL path
			// Expected format: /api/projects/{projectID}/...
			parts := strings.Split(request.URL.Path, "/")
			var projectIDStr string
			for i, part := range parts {
				if part == "projects" && i+1 < len(parts) {
					projectIDStr = parts[i+1]

					break
				}
			}

			if projectIDStr == "" || projectIDStr == "add" || projectIDStr == RecentKeyword {
				next.ServeHTTP(writer, request)

				return
			}

			projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
			if err != nil {
				errBadRequest := generatedapi.ErrorBadRequest{Error: generatedapi.ErrorBadRequestError{
					Message: generatedapi.NewOptString("invalid project ID"),
				}}
				errBadRequestData, _ := errBadRequest.MarshalJSON()
				http.Error(writer, string(errBadRequestData), http.StatusBadRequest)
				writer.Header().Set("Content-Type", "application/json; charset=utf-8")

				return
			}

			err = permissionsService.CanAccessProject(request.Context(), domain.ProjectID(projectID))
			if err != nil {
				slog.Error("failed to check project access", "error", err, "projectID", projectID)

				switch {
				case errors.Is(err, domain.ErrEntityNotFound):
					errNotFound := generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
						Message: generatedapi.NewOptString(err.Error()),
					}}
					errNotFoundData, _ := errNotFound.MarshalJSON()
					http.Error(writer, string(errNotFoundData), http.StatusNotFound)
					writer.Header().Set("Content-Type", "application/json; charset=utf-8")

					return
				case errors.Is(err, domain.ErrPermissionDenied):
					errPermDenied := generatedapi.ErrorPermissionDenied{Error: generatedapi.ErrorPermissionDeniedError{
						Message: generatedapi.NewOptString("permission denied"),
					}}
					errPermDeniedData, _ := errPermDenied.MarshalJSON()
					http.Error(writer, string(errPermDeniedData), http.StatusForbidden)
					writer.Header().Set("Content-Type", "application/json; charset=utf-8")

					return
				case errors.Is(err, domain.ErrUserNotFound):
					errUnauthorized := generatedapi.ErrorUnauthorized{Error: generatedapi.ErrorUnauthorizedError{
						Message: generatedapi.NewOptString("unauthorized"),
					}}
					errUnauthorizedData, _ := errUnauthorized.MarshalJSON()
					http.Error(writer, string(errUnauthorizedData), http.StatusUnauthorized)
					writer.Header().Set("Content-Type", "application/json; charset=utf-8")

					return
				default:
					errInternal := generatedapi.ErrorInternalServerError{
						Error: generatedapi.ErrorInternalServerErrorError{
							Message: generatedapi.NewOptString("internal server error"),
						},
					}
					errInternalData, _ := errInternal.MarshalJSON()
					http.Error(writer, string(errInternalData), http.StatusInternalServerError)
					writer.Header().Set("Content-Type", "application/json; charset=utf-8")

					return
				}
			}

			// Store project ID in context for later use
			ctx := wardencontext.WithProjectID(request.Context(), domain.ProjectID(projectID))
			next.ServeHTTP(writer, request.WithContext(ctx))
		})
	}
}

// ProjectManagement middleware checks if the user can manage the project.
//
//nolint:gocyclo // need refactoring
func ProjectManagement(permissionsService contract.PermissionsService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if request.Method == http.MethodGet {
				next.ServeHTTP(writer, request)

				return
			}

			// Extract project ID from URL path
			// Expected format: /api/projects/{projectID}/...
			parts := strings.Split(request.URL.Path, "/")
			var projectIDStr string
			for i, part := range parts {
				if part == "projects" && i+1 < len(parts) {
					projectIDStr = parts[i+1]

					break
				}
			}

			isStats := len(parts) >= 6 && parts[5] == "stats"

			if projectIDStr == "" || projectIDStr == "add" || projectIDStr == RecentKeyword || isStats {
				next.ServeHTTP(writer, request)

				return
			}

			isIssueManagement := len(parts) >= 7 && parts[5] == issuesStr && parts[7] == "change-status"

			projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
			if err != nil {
				errBadRequest := generatedapi.ErrorBadRequest{Error: generatedapi.ErrorBadRequestError{
					Message: generatedapi.NewOptString("invalid project ID"),
				}}
				errBadRequestData, _ := errBadRequest.MarshalJSON()
				http.Error(writer, string(errBadRequestData), http.StatusBadRequest)
				writer.Header().Set("Content-Type", "application/json; charset=utf-8")

				return
			}

			err = permissionsService.CanManageProject(
				request.Context(),
				domain.ProjectID(projectID),
				isIssueManagement,
			)
			if err != nil {
				slog.Error("failed to check project management permission", "error", err)

				switch {
				case errors.Is(err, domain.ErrEntityNotFound):
					errNotFound := generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
						Message: generatedapi.NewOptString(err.Error()),
					}}
					errNotFoundData, _ := errNotFound.MarshalJSON()
					http.Error(writer, string(errNotFoundData), http.StatusNotFound)
					writer.Header().Set("Content-Type", "application/json; charset=utf-8")

					return
				case errors.Is(err, domain.ErrPermissionDenied):
					errPermDenied := generatedapi.ErrorPermissionDenied{Error: generatedapi.ErrorPermissionDeniedError{
						Message: generatedapi.NewOptString("permission denied"),
					}}
					errPermDeniedData, _ := errPermDenied.MarshalJSON()
					http.Error(writer, string(errPermDeniedData), http.StatusForbidden)
					writer.Header().Set("Content-Type", "application/json; charset=utf-8")

					return
				case errors.Is(err, domain.ErrUserNotFound):
					errUnauthorized := generatedapi.ErrorUnauthorized{Error: generatedapi.ErrorUnauthorizedError{
						Message: generatedapi.NewOptString("unauthorized"),
					}}
					errUnauthorizedData, _ := errUnauthorized.MarshalJSON()
					http.Error(writer, string(errUnauthorizedData), http.StatusUnauthorized)
					writer.Header().Set("Content-Type", "application/json; charset=utf-8")

					return
				default:
					errInternal := generatedapi.ErrorInternalServerError{
						Error: generatedapi.ErrorInternalServerErrorError{
							Message: generatedapi.NewOptString("internal server error"),
						},
					}
					errInternalData, _ := errInternal.MarshalJSON()
					http.Error(writer, string(errInternalData), http.StatusInternalServerError)
					writer.Header().Set("Content-Type", "application/json; charset=utf-8")

					return
				}
			}

			// Store project ID in context for later use
			ctx := wardencontext.WithProjectID(request.Context(), domain.ProjectID(projectID))
			next.ServeHTTP(writer, request.WithContext(ctx))
		})
	}
}

// IssueAccess middleware checks if the user has access to the issue.
func IssueAccess(permissionsService contract.PermissionsService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			// Extract issue ID from URL path
			// Expected format: /api/projects/{projectID}/issues/{issueID}/...
			parts := strings.Split(request.URL.Path, "/")
			var issueIDStr string
			for i, part := range parts {
				if part == issuesStr && i+1 < len(parts) {
					issueIDStr = parts[i+1]

					break
				}
			}

			if issueIDStr == "" || issueIDStr == RecentKeyword || issueIDStr == "timeseries" {
				next.ServeHTTP(writer, request)

				return
			}

			issueID, err := strconv.ParseUint(issueIDStr, 10, 64)
			if err != nil {
				errBadRequest := generatedapi.ErrorBadRequest{Error: generatedapi.ErrorBadRequestError{
					Message: generatedapi.NewOptString("invalid issue ID"),
				}}
				errBadRequestData, _ := errBadRequest.MarshalJSON()
				http.Error(writer, string(errBadRequestData), http.StatusBadRequest)
				writer.Header().Set("Content-Type", "application/json; charset=utf-8")

				return
			}

			err = permissionsService.CanAccessIssue(request.Context(), domain.IssueID(issueID))
			if err != nil {
				slog.Error("failed to check issue access", "error", err)

				switch {
				case errors.Is(err, domain.ErrEntityNotFound):
					errNotFound := generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
						Message: generatedapi.NewOptString(err.Error()),
					}}
					errNotFoundData, _ := errNotFound.MarshalJSON()
					http.Error(writer, string(errNotFoundData), http.StatusNotFound)
					writer.Header().Set("Content-Type", "application/json; charset=utf-8")

					return
				case errors.Is(err, domain.ErrPermissionDenied):
					errPermDenied := generatedapi.ErrorPermissionDenied{Error: generatedapi.ErrorPermissionDeniedError{
						Message: generatedapi.NewOptString("permission denied"),
					}}
					errPermDeniedData, _ := errPermDenied.MarshalJSON()
					http.Error(writer, string(errPermDeniedData), http.StatusForbidden)
					writer.Header().Set("Content-Type", "application/json; charset=utf-8")

					return
				case errors.Is(err, domain.ErrUserNotFound):
					errUnauthorized := generatedapi.ErrorUnauthorized{Error: generatedapi.ErrorUnauthorizedError{
						Message: generatedapi.NewOptString("unauthorized"),
					}}
					errUnauthorizedData, _ := errUnauthorized.MarshalJSON()
					http.Error(writer, string(errUnauthorizedData), http.StatusUnauthorized)
					writer.Header().Set("Content-Type", "application/json; charset=utf-8")

					return
				default:
					errInternal := generatedapi.ErrorInternalServerError{
						Error: generatedapi.ErrorInternalServerErrorError{
							Message: generatedapi.NewOptString("internal server error"),
						},
					}
					errInternalData, _ := errInternal.MarshalJSON()
					http.Error(writer, string(errInternalData), http.StatusInternalServerError)
					writer.Header().Set("Content-Type", "application/json; charset=utf-8")

					return
				}
			}

			next.ServeHTTP(writer, request)
		})
	}
}

// IssueManagement middleware checks if the user can manage the issue.
func IssueManagement(permissionsService contract.PermissionsService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if request.Method == http.MethodGet {
				next.ServeHTTP(writer, request)

				return
			}

			// Extract issue ID from URL path
			// Expected format: /api/projects/{projectID}/issues/{issueID}/...
			parts := strings.Split(request.URL.Path, "/")
			var issueIDStr string
			for i, part := range parts {
				if part == issuesStr && i+1 < len(parts) {
					issueIDStr = parts[i+1]

					break
				}
			}

			if issueIDStr == "" ||
				issueIDStr == RecentKeyword ||
				issueIDStr == "timeseries" {
				next.ServeHTTP(writer, request)

				return
			}

			issueID, err := strconv.ParseUint(issueIDStr, 10, 64)
			if err != nil {
				errBadRequest := generatedapi.ErrorBadRequest{Error: generatedapi.ErrorBadRequestError{
					Message: generatedapi.NewOptString("invalid issue ID"),
				}}
				errBadRequestData, _ := errBadRequest.MarshalJSON()
				http.Error(writer, string(errBadRequestData), http.StatusBadRequest)
				writer.Header().Set("Content-Type", "application/json; charset=utf-8")

				return
			}

			err = permissionsService.CanManageIssue(request.Context(), domain.IssueID(issueID))
			if err != nil {
				slog.Error("failed to check issue management permission", "error", err)

				switch {
				case errors.Is(err, domain.ErrEntityNotFound):
					errNotFound := generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
						Message: generatedapi.NewOptString(err.Error()),
					}}
					errNotFoundData, _ := errNotFound.MarshalJSON()
					http.Error(writer, string(errNotFoundData), http.StatusNotFound)
					writer.Header().Set("Content-Type", "application/json; charset=utf-8")

					return
				case errors.Is(err, domain.ErrPermissionDenied):
					errPermDenied := generatedapi.ErrorPermissionDenied{Error: generatedapi.ErrorPermissionDeniedError{
						Message: generatedapi.NewOptString("permission denied"),
					}}
					errPermDeniedData, _ := errPermDenied.MarshalJSON()
					http.Error(writer, string(errPermDeniedData), http.StatusForbidden)
					writer.Header().Set("Content-Type", "application/json; charset=utf-8")

					return
				case errors.Is(err, domain.ErrUserNotFound):
					errUnauthorized := generatedapi.ErrorUnauthorized{Error: generatedapi.ErrorUnauthorizedError{
						Message: generatedapi.NewOptString("unauthorized"),
					}}
					errUnauthorizedData, _ := errUnauthorized.MarshalJSON()
					http.Error(writer, string(errUnauthorizedData), http.StatusUnauthorized)
					writer.Header().Set("Content-Type", "application/json; charset=utf-8")

					return
				default:
					errInternal := generatedapi.ErrorInternalServerError{
						Error: generatedapi.ErrorInternalServerErrorError{
							Message: generatedapi.NewOptString("internal server error"),
						},
					}
					errInternalData, _ := errInternal.MarshalJSON()
					http.Error(writer, string(errInternalData), http.StatusInternalServerError)
					writer.Header().Set("Content-Type", "application/json; charset=utf-8")

					return
				}
			}

			next.ServeHTTP(writer, request)
		})
	}
}
