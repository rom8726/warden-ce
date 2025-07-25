package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	wardencontext "github.com/rom8726/warden/internal/context"
	generatedapi "github.com/rom8726/warden/internal/generated/ingestserver"
	"github.com/rom8726/warden/pkg/metrics"
)

func (r *RestAPI) StoreEvent(
	ctx context.Context,
	req generatedapi.StoreEventReq,
	_ generatedapi.StoreEventParams,
) (*generatedapi.StoreEventResponse, error) {
	// Convert StoreEventReq to map[string]any
	eventData := make(map[string]any)
	for k, v := range req {
		var value any
		if err := json.Unmarshal(v, &value); err != nil {
			slog.Error("Failed to unmarshal event data", "error", err, "key", k)
			metrics.ValidationErrors.WithLabelValues("invalid_json").Inc()

			return nil, fmt.Errorf("unmarshal event data: %w", err)
		}
		eventData[k] = value
	}

	eventID, err := r.storeEventUseCase.StoreEvent(ctx, wardencontext.ProjectID(ctx), eventData)
	if err != nil {
		slog.Error("store event failed", "error", err)

		return nil, err
	}

	return &generatedapi.StoreEventResponse{
		ID: string(eventID),
	}, nil
}
