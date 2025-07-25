package rest

import (
	"context"
	"log/slog"
	"time"

	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

// GetVersions implements the GetVersions operation.
func (r *RestAPI) GetVersions(ctx context.Context) (generatedapi.GetVersionsRes, error) {
	// Get versions from all components
	components, err := r.versionsUseCase.GetVersions(ctx)
	if err != nil {
		slog.Error("failed to get versions", "error", err)

		return &generatedapi.ErrorInternalServerError{
			Error: generatedapi.ErrorInternalServerErrorError{
				Message: generatedapi.NewOptString("failed to get versions"),
			},
		}, nil
	}

	// Convert to generated API types
	apiComponents := make([]generatedapi.ComponentVersion, 0, len(components))
	for _, comp := range components {
		status := generatedapi.ComponentVersionStatusAvailable
		if comp.Status == "unavailable" {
			status = generatedapi.ComponentVersionStatusUnavailable
		}

		apiComponents = append(apiComponents, generatedapi.ComponentVersion{
			Name:      comp.Name,
			Version:   comp.Version,
			BuildTime: comp.BuildTime,
			Status:    status,
		})
	}

	return &generatedapi.VersionsResponse{
		Components:  apiComponents,
		CollectedAt: time.Now(),
	}, nil
}
