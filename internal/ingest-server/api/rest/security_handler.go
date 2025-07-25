package rest

import (
	"context"
	"errors"
	"strings"

	wardencontext "github.com/rom8726/warden/internal/context"
	generatedapi "github.com/rom8726/warden/internal/generated/ingestserver"
	"github.com/rom8726/warden/internal/ingest-server/contract"
)

var _ generatedapi.SecurityHandler = (*SecurityHandler)(nil)

type SecurityHandler struct {
	projectService contract.ProjectsUseCase
}

func NewSecurityHandler(
	projectService contract.ProjectsUseCase,
) *SecurityHandler {
	return &SecurityHandler{
		projectService: projectService,
	}
}

func (r *SecurityHandler) HandleSentryAuth(
	ctx context.Context,
	_ generatedapi.OperationName,
	tokenHolder generatedapi.SentryAuth,
) (context.Context, error) {
	sentryKey := parseSentryKeyFromAuth(tokenHolder.APIKey)
	projectID := wardencontext.ProjectID(ctx)

	allowed, err := r.projectService.ValidateProjectKey(ctx, projectID, sentryKey)
	if err != nil {
		return nil, err
	}

	if !allowed {
		return nil, errors.New("invalid or unauthorized key")
	}

	return ctx, nil
}

func parseSentryKeyFromAuth(auth string) string {
	parts := strings.Split(auth, ",")
	for _, part := range parts {
		if strings.Contains(part, "sentry_key=") {
			kv := strings.Split(strings.TrimSpace(part), "=")
			if len(kv) >= 2 {
				return strings.TrimSpace(strings.Trim(kv[1], `"`))
			}
		}
	}

	return ""
}
