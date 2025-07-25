package rest

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	wardencontext "github.com/rom8726/warden/internal/context"
	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/ingestserver"
)

func (r *RestAPI) ReceiveEnvelope(
	ctx context.Context,
	req generatedapi.ReceiveEnvelopeReq,
	_ generatedapi.ReceiveEnvelopeParams,
) (generatedapi.ReceiveEnvelopeRes, error) {
	err := r.envelopeUseCase.ReceiveEnvelope(ctx, wardencontext.ProjectID(ctx), req.Data)
	if err != nil {
		slog.Error("receive envelope failed", "error", err)

		switch {
		case errors.Is(err, domain.ErrNoEnvelope), errors.Is(err, domain.ErrInvalidEnvelopeHeader):
			return &generatedapi.ErrorBadRequest{Error: generatedapi.ErrorBadRequestError{
				Message: generatedapi.NewOptString(err.Error()),
			}}, nil
		default:
			return &generatedapi.ReceiveEnvelopeOK{}, err
		}
	}

	return &generatedapi.ReceiveEnvelopeOK{
		Data: strings.NewReader("OK"),
	}, nil
}
