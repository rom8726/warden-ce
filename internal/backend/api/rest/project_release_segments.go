package rest

import (
	"context"
	"errors"

	"github.com/rom8726/warden/internal/backend/dto"
	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

const unknownKeyword = "unknown"

func (r *RestAPI) GetProjectReleaseSegments(
	ctx context.Context,
	params generatedapi.GetProjectReleaseSegmentsParams,
) (generatedapi.GetProjectReleaseSegmentsRes, error) {
	projectID := domain.ProjectID(params.ProjectID)
	release := params.Release
	segment := domain.SegmentName(params.Segment)

	segments, err := r.analyticsUseCase.GetUserSegments(ctx, projectID, release)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrEntityNotFound):
			return nil, &generatedapi.ErrorStatusCode{StatusCode: 404}
		default:
			return nil, &generatedapi.ErrorStatusCode{StatusCode: 500}
		}
	}

	var values map[string]uint
	switch segment {
	case domain.SegmentNamePlatform:
		values = map[string]uint{}
		for k, v := range segments.Platform {
			values[stringKey(k)] = v
		}
	case domain.SegmentNameBrowserName:
		values = map[string]uint{}
		for k, v := range segments.Browser {
			values[stringKey(k)] = v
		}
	case domain.SegmentNameOSName:
		values = map[string]uint{}
		for k, v := range segments.OS {
			values[stringKey(k)] = v
		}
	case domain.SegmentNameDeviceArch:
		values = map[string]uint{}
		for k, v := range segments.DeviceArch {
			values[stringKey(k)] = v
		}
	case domain.SegmentNameRuntimeName:
		values = map[string]uint{}
		for k, v := range segments.RuntimeName {
			values[stringKey(k)] = v
		}
	default:
		return nil, &generatedapi.ErrorStatusCode{StatusCode: 400}
	}

	resp := dto.ToReleaseSegmentsResponse(string(segment), values)

	return &resp, nil
}

func stringKey(key domain.UserSegmentKey) string {
	if key == "" {
		return unknownKeyword
	}

	return string(key)
}
